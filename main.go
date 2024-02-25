// expulo extract purge and lod data in two PostgreSQL instances
package main

import (
	"fmt"
	"flag"
	"os"
	"strings"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func fullname(schemaname string, datname string, attname string) string {
	return fmt.Sprintf("%s.%s.%s", schemaname, datname, attname)
}

func fullTableName(schemaname string, datname string) string {
	return fmt.Sprintf("%s.%s", schemaname, datname)
}


func init() {
	log.SetLevel(log.DebugLevel)
	//	log.SetLevel(log.InfoLevel)

}

func main() {
	version := "0.0.2"
	// Command line flag
	var purgeOnly bool
	flag.BoolVar(&purgeOnly, "purge", false, "Only purge destination tables and exit, no data will be inserted")

	flag.Parse()

	// Read connection information from env variables
	src_host := os.Getenv("PGSRCHOST")
	src_port := os.Getenv("PGSRCPORT")
	src_user := os.Getenv("PGSRCUSER")
	src_pass := os.Getenv("PGSRCPASSWORD")
	src_db   := os.Getenv("PGSRCDATABASE")

	dst_host := os.Getenv("PGDSTHOST")
	dst_port := os.Getenv("PGDSTPORT")
	dst_user := os.Getenv("PGDSTUSER")
	dst_pass := os.Getenv("PGDSTPASSWORD")
	dst_db   := os.Getenv("PGDSTDATABASE")

	// Construct connection string
	conxSource, dsn_src := get_dsn(src_host, src_port, src_user, src_pass, src_db, version)
	conxDestination, dsn_dst := get_dsn(dst_host, dst_port, dst_user, dst_pass, dst_db, version)

	// Connect to the database source
	log.Debug("Connect on source")
	db_src := connectDb(conxSource)
	log.Info(fmt.Sprintf("Use %s as source", dsn_src))

	// Connect to the database destination
	log.Debug("Connect on destination")
	db_dst := connectDb(conxDestination)
	log.Info(fmt.Sprintf("Use %s as destination", dsn_dst))

	// Read the configuration
	config := read_config("config.json")
	log.Debug("Read config done")
	log.Debug("Number of tables found in conf: ", len(config.Tables))

	// Delete data on destination tables
	purge_destination(config, db_dst)


	// if command line parameter set do purge and exit
	if purgeOnly == true {
		log.Debug("Exit on option, purge")
		os.Exit(0)
	}

	for _, t := range config.Tables {
		table_name := t.Name
		batch_size := 4
		src_query := fmt.Sprintf("SELECT * FROM %s WHERE id >= $1 AND id < $2", table_name)

		// Filter the data on source to fetch a subset of rows in a table
		if len(t.Filter) > 0 {
			src_query = fmt.Sprintf("%s AND %s", src_query, t.Filter)
		}
		log.Debug(src_query)

		keepRunning := true
		run := 0
		for keepRunning {
			rows, err := db_src.Query(src_query, batch_size * run, batch_size * (run + 1))
			run = run + 1
			if err != nil {
				fmt.Println("Error executing query:", err)
				os.Exit(1)
			}
			defer rows.Close()

			columns, err := rows.Columns()
			if err != nil {
				fmt.Println("Error getting column names:", err)
				return
			}

			count := 0
			for rows.Next() {
				var colnames []string
				count = count + 1
				cols := make([]interface{}, len(columns))

				columnPointers := make([]interface{}, len(cols))

				for i, _ := range cols {
					columnPointers[i] = &cols[i]

				}
				rows.Scan(columnPointers...)
				nbcol := 1
				var colparam []string
				var colvalue []interface{}
				//fval := make([]interface{}, len(cols))
				// Manage what we do it data here
				for i, _ := range cols {
					cfvalue := "notfound"
					col, err := get_cols(t, columns[i])
					if err == 1 {
						cfvalue = "notfound"
					} else {
						cfvalue = col.Generator
					}

					//colname := fullname(t.Schema, t.Name, columns[i])
					//cfvalue := config[colname]

					//log.Output(1, fmt.Sprintf("%s %s", colname, cfvalue))

					// If the configuration ignore the column it won't be present
					// in the INSERT statement
					if cfvalue != "ignore" {

						colnames = append(colnames, columns[i])

						// Assign the target value
						switch cfvalue{
						case "null":
							colvalue = append(colvalue, nil)
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "mask":
							colvalue = append(colvalue, mask())
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "randomInt":
							colvalue = append(colvalue, randomInt())
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "randomIntMinMax":
							colvalue = append(colvalue, randomIntMinMax(col.Min, col.Max))
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "randomFloat":
							colvalue = append(colvalue, randomFloat())
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "randomString":
							colvalue = append(colvalue, randomString())
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "md5":
							colvalue = append(colvalue, md5signature(fmt.Sprintf("%v", cols[i])))
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "randomTimeTZ":
							colvalue = append(colvalue, randomTimeTZ(col.Timezone))
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						case "sql":
							nbcol = nbcol - 1
							colparam = append(colparam, col.SQLFunction)
						default:
							colvalue = append(colvalue, cols[i])
							colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
						}

						nbcol = nbcol +1

					}
				}

				col_names := strings.Join(colnames, ",")

				dst_query := "INSERT INTO " + table_name + " (" + col_names + ") VALUES ("+strings.Join(colparam,",") + ")"
				log.Debug(dst_query)
				_, err := db_dst.Exec(dst_query, colvalue...)
				if err != nil {
					log.Error("Error during INSERT on :", table_name, err)
					return
				}

			}
			if count == 0 { keepRunning = false }
			log.Info(fmt.Sprintf("%d",count))
		}


	}
}
