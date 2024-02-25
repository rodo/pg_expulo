// expulo extract purge and lod data in two PostgreSQL instances
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	//	log.SetLevel(log.InfoLevel)

}

func main() {
	version := "0.0.2"
	// Command line flag
	var purgeOnly bool
	var tryOnly bool
	flag.BoolVar(&purgeOnly, "purge", false, "Only purge destination tables and exit, no data will be inserted")
	flag.BoolVar(&tryOnly, "try", false, "ROLLBACK everything on target. No data will be inserted")

	flag.Parse()

	// Read connection information from env variables
	srcHost := os.Getenv("PGSRCHOST")
	srcPort := os.Getenv("PGSRCPORT")
	srcUser := os.Getenv("PGSRCUSER")
	srcPass := os.Getenv("PGSRCPASSWORD")
	srcDb := os.Getenv("PGSRCDATABASE")

	dstHost := os.Getenv("PGDSTHOST")
	dstPort := os.Getenv("PGDSTPORT")
	dstUser := os.Getenv("PGDSTUSER")
	dstPass := os.Getenv("PGDSTPASSWORD")
	dstDb := os.Getenv("PGDSTDATABASE")

	// Construct connection string
	conxSource, dsnSrc := get_dsn(srcHost, srcPort, srcUser, srcPass, srcDb, version)
	conxDestination, dsnDst := get_dsn(dstHost, dstPort, dstUser, dstPass, dstDb, version)

	// Connect to the database source
	log.Debug("Connect on source")
	dbSrc := connectDb(conxSource)
	log.Info(fmt.Sprintf("Use %s as source", dsnSrc))

	// Connect to the database destination
	log.Debug("Connect on destination")
	dbDst := connectDb(conxDestination)
	log.Info(fmt.Sprintf("Use %s as destination", dsnDst))

	// Start a transaction
	txDst, err := dbDst.Begin()
	if err != nil {
		log.Fatal("Error starting transaction: ", err)
	}
	defer txDst.Rollback() // Rollback the transaction if it hasn't been committed

	// Read the configuration
	config := read_config("config.json")
	log.Debug("Read config done")
	log.Debug("Number of tables found in conf: ", len(config.Tables))

	// Delete data on destination tables
	purge_destination(config, dbDst)

	// if command line parameter set do purge and exit
	if purgeOnly == true {
		log.Debug("Exit on option, purge")
		os.Exit(0)
	}

	for _, t := range config.Tables {
		tableFullname := fullTableName(t.Schema, t.Name)
		//batch_size := 4
		src_query := fmt.Sprintf("SELECT * FROM %s WHERE true", tableFullname)

		// Filter the data on source to fetch a subset of rows in a table
		if len(t.Filter) > 0 {
			src_query = fmt.Sprintf("%s AND %s", src_query, t.Filter)
		}
		doTables(dbSrc, txDst, t, src_query)
	}

	if tryOnly {
		// Rollback the transaction on target as requested
		if err := txDst.Rollback(); err != nil {
			log.Fatal("Error committing transaction: ", err)
		}
		log.Info("Rollback on target")
	} else {
		// Commit the transaction on target if all queries succeed
		if err := txDst.Commit(); err != nil {
			log.Fatal("Error committing transaction: ", err)
		} else {
			log.Info("Commit on target")
		}
	}

}

func doTables(dbSrc *sql.DB, txDst *sql.Tx, t Table, src_query string) {
	tableFullname := fullTableName(t.Schema, t.Name)

	rows, err := dbSrc.Query(src_query)
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
	var multirows [][]interface{}

	count := 0
	nbinsert := 0
	var colnames []string
	var colparam []string
	for rows.Next() {
		colnames = []string{}
		count = count + 1
		nbinsert = nbinsert + 1
		cols := make([]interface{}, len(columns))

		columnPointers := make([]interface{}, len(cols))

		for i, _ := range cols {
			columnPointers[i] = &cols[i]

		}
		rows.Scan(columnPointers...)
		nbcol := 1
		colparam = []string{}
		var colvalue []interface{}
		//fval := make([]interface{}, len(cols))
		// Manage what we do it data here
		for i, _ := range cols {
			cfvalue := "notfound"
			col, found := get_cols(t, columns[i])
			if found {
				cfvalue = col.Generator
			} else {
				cfvalue = "notfound"
			}

			// If the configuration ignore the column it won't be present
			// in the INSERT statement
			if cfvalue != "ignore" {

				colnames = append(colnames, columns[i])

				// Assign the target value
				switch cfvalue {
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
				nbcol = nbcol + 1
			}
		}

		// INSERT
		multirows = append(multirows, colvalue)

		if nbinsert > 9 {
			log.Debug(fmt.Sprintf("Insert %d rows in table ", nbinsert))
			nbinsert = 0
			insertMultiData(txDst, tableFullname, colnames, colparam, multirows)
			multirows = multirows[:0]
		}
	}
	insertMultiData(txDst, tableFullname, colnames, colparam, multirows)
	log.Debug(fmt.Sprintf("Inserted %d rows in table %s", count, t.Name))

}

func insertMultiData(dbDst *sql.Tx, tableFullname string, colnames []string, colparam []string, multirows [][]interface{}) {
	col_names := strings.Join(colnames, ",")

	nbColumns := len(colnames)
	nbRows := len(multirows)

	pat := toolPat(nbRows, nbColumns, colparam)

	//log.Debug(fmt.Sprintf("there is %d rows of %d columns", nbRows, nbColumns))

	var allValues []interface{}
	for _, row := range multirows {
		// Append each element of the row to allValues
		allValues = append(allValues, row...)
	}

	destQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableFullname, col_names, pat)

	_, err := dbDst.Exec(destQuery, allValues...)
	if err != nil {
		log.Debug("Error during INSERT on :", err)
		log.Debug(destQuery)
		log.Debug(allValues)
		log.Fatal("Error during INSERT on :", tableFullname)

		return
	}
}
