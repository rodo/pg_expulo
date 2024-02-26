// expulo extract purge and lod data in two PostgreSQL instances
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// Config store the whole configuration read from json file
type Config struct {
	Tables []Table `json:"tables"`
}

// Columns contains a collection of Column
type Columns struct {
	Columns []Column `json:"columns"`
}

// Table represent a table with her property in configuration file
type Table struct {
	Name           string   `json:"name"`
	Columns        []Column `json:"columns"`
	Schema         string   `json:"schema"`
	CleanMethod    string   `json:"clean"`
	Filter         string   `json:"filter"`
	DeletionFilter string   `json:"deletion_filter"`
}

type Column struct {
	Name        string `json:"name"`
	Generator   string `json:"generator"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Timezone    string `json:"timezone"`
	SQLFunction string `json:"function"`
}

var (
	Version = "0.0.2"
)

func init() {
	log.SetLevel(log.DebugLevel)
	// log.SetLevel(log.InfoLevel)

}

func CheckFlags() (bool, bool, string) {

	version := flag.Bool("version", false, "display version")
	tryOnly := flag.Bool("try", false, "ROLLBACK everything on target. No data will be inserted")
	purgeOnly := flag.Bool("purge", false, "Only purge destination tables and exit, no data will be inserted")
	configFile := flag.String("config", "config.json", "Configuration file to use")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	return *purgeOnly, *tryOnly, *configFile
}

func main() {

	// Command line flag
	purgeOnly, tryOnly, configFile := CheckFlags()

	// Read the configuration
	config := readConfig(configFile)
	log.Debug("Read config done")
	log.Debug("Number of tables found in conf: ", len(config.Tables))

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
	conxSource, dsnSrc := getDsn(srcHost, srcPort, srcUser, srcPass, srcDb, Version)
	conxDestination, dsnDst := getDsn(dstHost, dstPort, dstUser, dstPass, dstDb, Version)

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

	// Delete data on destination tables
	purgeTarget(config, txDst)

	// if command line parameter set do purge and exit
	if purgeOnly == true {
		log.Debug("Exit on option, purge")

	} else {

		// Loop over all tables configured
		for _, t := range config.Tables {
			tableFullname := fullTableName(t.Schema, t.Name)

			src_query := fmt.Sprintf("SELECT * FROM %s WHERE true", tableFullname)

			// Filter the data on source to fetch a subset of rows in a table
			if len(t.Filter) > 0 {
				src_query = fmt.Sprintf("%s AND %s", src_query, t.Filter)
			}
			startTime := time.Now()
			nbrows := doTable(dbSrc, txDst, t, src_query)
			elapsedTime := time.Since(startTime)
			log.Info(fmt.Sprintf("%s : inserted %d rows total in %s", tableFullname, nbrows, elapsedTime))
		}
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

func doTable(dbSrc *sql.DB, txDst *sql.Tx, t Table, src_query string) int {
	tableFullname := fullTableName(t.Schema, t.Name)

	log.Info(fmt.Sprintf("%s : read data fom table", tableFullname))

	rows, columns := queryTableSource(dbSrc, src_query)

	var multirows [][]interface{}
	lenColumns := len(columns)

	count := 0
	nbinsert := 0
	var colnames []string
	var colparam []string

	for rows.Next() {
		colnames = []string{}
		count++
		nbinsert++
		cols := make([]interface{}, lenColumns)

		columnPointers := make([]interface{}, len(cols))

		for i, _ := range cols {
			columnPointers[i] = &cols[i]

		}
		rows.Scan(columnPointers...)
		nbcol := 1
		colparam = []string{}
		var colvalue []interface{}

		// Manage what we do it data here
		for i, _ := range cols {
			cfvalue := "notfound"
			col, found := getCols(t, columns[i])
			if found {
				cfvalue = col.Generator
			} else {
				cfvalue = "notfound"
			}

			// If the configuration ignore the column it won't be present
			// in the INSERT statement

			colvalue, colparam, nbcol, colnames = fillColumn(col, cfvalue, colvalue, colparam, nbcol, cols, colnames, i, columns)

		}

		// INSERT
		multirows = append(multirows, colvalue)

		batch_size := 1000
		if nbinsert > batch_size-1 {
			log.Debug(fmt.Sprintf("Insert %d rows in table %s", nbinsert, t.Name))
			nbinsert = 0
			insertMultiData(txDst, tableFullname, colnames, colparam, multirows)
			multirows = multirows[:0]
		}
	}
	insertMultiData(txDst, tableFullname, colnames, colparam, multirows)
	return count
}
