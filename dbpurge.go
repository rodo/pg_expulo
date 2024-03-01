package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:embed sql/linked_tables.sql
var qryLinkedTables string

// Purge all the tables in the database target

func getTableByName(config Config, name string) (Table, bool) {
	for _, table := range config.Tables {
		if table.Name == name {
			return table, true
		}
	}
	return Table{}, false
}

func getTableByFullName(config Config, name string) (Table, bool) {
	for _, table := range config.Tables {
		if fullTableName(table.Schema, table.Name) == name {
			return table, true
		}
	}
	return Table{}, false
}

func purgeTarget(config Config, txDst *sql.Tx) {

	forcePurge := true
	var tableList []string
	var OrderedTables []Table
	// Order tables depending on foreign keys
	for _, t := range config.Tables {
		tableList = append(tableList, fullTableName(t.Schema, t.Name))
	}

	for _, tname := range getOrderTableList(tableList, txDst) {
		// tableList = append(tableList, fmt.Sprintf("%s.%s", t.Schema, t.Name))
		t, found := getTableByFullName(config, tname)
		if found {
			OrderedTables = append(OrderedTables, t)
		}
	}

	// List all tables in purge order
	for _, t := range OrderedTables {
		log.Debug(fmt.Sprintf("Will clean table : %s.%s with %s", t.Schema, t.Name, t.CleanMethod))
	}

	// Loop over all tables found in configuration file
	for _, t := range OrderedTables {
		tableFullname := fullTableName(t.Schema, t.Name)

		log.Info(fmt.Sprintf("%s : clean table (method:%s, filter:%s)", tableFullname, t.CleanMethod, t.Filter))

		// Clean target tables
		switch t.CleanMethod {
		case "append":
			log.Debug("Do nothing on target purge according to configuration")

		case "truncate":
			log.Debug("TRUNCATE TABLE according to default")
			dstQuery := "TRUNCATE " + tableFullname + ";"
			_, err := txDst.Exec(dstQuery)
			if err != nil {
				if forcePurge {
					log.Error(err)
				} else {
					log.Fatal(err)
				}
			}
		default:
			_ = deleteData(t, forcePurge, txDst)
		}
	}
}

func deleteData(t Table, forcePurge bool, txDst *sql.Tx) error {
	log.Debug(fmt.Sprintf("DELETE data from %s according to configuration", t.Name))

	var dstQry string

	if len(t.DeletionFilter) > 0 {
		dstQry = fmt.Sprintf("DELETE FROM %s.%s WHERE %s", t.Schema, t.Name, t.DeletionFilter)
	} else {
		dstQry = fmt.Sprintf("DELETE FROM %s.%s", t.Schema, t.Name)
	}

	log.Debug(dstQry)

	_, err := txDst.Exec(dstQry)
	if err != nil {
		if forcePurge {
			log.Error(err)
		} else {
			log.Fatal(err)
		}

		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			// Print the error code and message
			log.Printf("PostgreSQL error code: %s\n", pqErr.Code)
			log.Printf("PostgreSQL error message: %s\n", pqErr.Message)
		} else {
			// For non-specific errors, just log the error message
			log.Println("Error inserting into tableB: ", err)
		}

	}
	return err
}

// Order the table list on number of foreign keys poiting to them
// This will ensure to purge in first table with no foriegn keys that
// pointing to them
// The order is not perfect as it is based on numer of foreign keys
// it's a first step
func getOrderTableList(tableList []string, txDst *sql.Tx) []string {

	var pkName string
	var nbfk int
	var orderedTableList []string
	tables := "{" + strings.Join(tableList, ",") + "}"

	// Query data from tableA
	rows, erri := txDst.Query(qryLinkedTables, tables)
	if erri != nil {
		log.Fatal("Error querying data from tableA: ", erri)
	}
	// Iterate through the rows from tableA and insert into tableB
	for rows.Next() {
		if erri := rows.Scan(&pkName, &nbfk); erri != nil {
			log.Error("Error scanning row: ", erri)

		}
		orderedTableList = append(orderedTableList, pkName)
	}
	rows.Close()
	log.Debug(pkName)
	return orderedTableList
}
