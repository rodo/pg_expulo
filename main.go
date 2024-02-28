// expulo extract purge and lod data in two PostgreSQL instances
package main

import (
	"database/sql"
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

// Table represent a table with her property in configuration file
type TriggerConstraint struct {
	TableFullName  string
	ConstraintName string
}

var (
	Version    = "0.0.2"
	tryOnly    = false
	purgeOnly  = false
	configFile = "config.json"
)

func init() {
	// Set default LogLevel to INFO
	log.SetLevel(log.InfoLevel)

	// Check if stdout is connected to a terminal
	// If not remove colors in logs to be friendly
	if !IsTerminal(os.Stdout) {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
}

func main() {
	// Parse command line arguments
	flagParse()

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

	// Build an Array with all table names in fullname form
	var tableList []string
	for _, t := range config.Tables {
		tableList = append(tableList, fullTableName(t.Schema, t.Name))
	}

	log.Debug("tableList contains : ", tableList)

	// Read the foreign keys
	triggerConstraints := GetTriggerConstraints(dbDst, tableList)

	// Delete data on destination tables
	DeferForeignKeys(dbDst, triggerConstraints)
	purgeTarget(config, txDst)

	// if command line parameter set do purge and exit
	if purgeOnly == true {
		log.Debug("Exit on option, purge")
		CloseTx(txDst, tryOnly)
		ReactivateForeignKeys(dbDst, triggerConstraints)
		os.Exit(0)
	}

	// List all tables in insert order
	for _, t := range config.Tables {
		tableFullname := fullTableName(t.Schema, t.Name)
		log.Debug(fmt.Sprintf("Will insert in : %s", tableFullname))
	}

	// Loop over all tables configured
	for _, t := range config.Tables {
		tableFullname := fullTableName(t.Schema, t.Name)

		src_query := fmt.Sprintf("SELECT * FROM %s WHERE true", tableFullname)

		// Filter the data on source to fetch a subset of rows in a table
		if len(t.Filter) > 0 {
			src_query = fmt.Sprintf("%s AND %s", src_query, t.Filter)
		}
		startTime := time.Now()
		nbrows, _ := doTable(dbSrc, txDst, t, src_query)
		elapsedTime := time.Since(startTime)
		log.Info(fmt.Sprintf("%s : inserted %d rows total in %s", tableFullname, nbrows, elapsedTime))

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
	log.Info("Thank you for using pg_expulo")
}

func doTable(dbSrc *sql.DB, txDst *sql.Tx, t Table, src_query string) (int, string) {
	tableFullname := fullTableName(t.Schema, t.Name)

	log.Info(fmt.Sprintf("%s : read data fom table", tableFullname))

	rows, columns := queryTableSource(dbSrc, src_query)

	var multirows [][]interface{}
	lenColumns := len(columns)

	count := 0
	nbinsert := 0
	var errCode string
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
			_, errCode = insertMultiData(txDst, tableFullname, colnames, colparam, multirows)
			multirows = multirows[:0]
		}
	}
	_, errCode = insertMultiData(txDst, tableFullname, colnames, colparam, multirows)
	return count, errCode
}

func RemoveAtIndex(slice []Table, index int) []Table {
	return append(slice[:index], slice[index+1:]...)
}
