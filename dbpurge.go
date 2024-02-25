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

//go:embed sql/primary_key.sql
var qry_primary_key string

//go:embed sql/linked_tables.sql
var qry_linked_tables string

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
	var table_list []string
	var OrderedTables []Table
	// Order table depending on foreign keys
	for _, t := range config.Tables {
		table_list = append(table_list, fullTableName(t.Schema, t.Name))
	}

	for _, tname := range OrderTableList(table_list, txDst) {
		// table_list = append(table_list, fmt.Sprintf("%s.%s", t.Schema, t.Name))
		t, found := getTableByFullName(config, tname)
		if found {
			OrderedTables = append(OrderedTables, t)
		}
	}

	for _, t := range OrderedTables {
		log.Debug(fmt.Sprintf("Will clean table : %s.%s with %s", t.Schema, t.Name, t.CleanMethod))
	}

	// Loop over all tables found in configuration file
	for _, t := range OrderedTables {
		table_name := t.Name

		log.Info(fmt.Sprintf("Clean table : %s (%s, %s)", t.Name, t.CleanMethod, t.Filter))

		// Clean target tables
		switch t.CleanMethod {
		case "append":
			log.Debug("Do nothing on target purge according to configuration")
		case "delete":
			_ = deleteData(t, forcePurge, txDst)
		default:
			log.Debug("TRUNCATE TABLE according to default")
			dst_query := "TRUNCATE " + table_name + ";"
			_, err := txDst.Exec(dst_query)
			if err != nil {
				if forcePurge {
					log.Error(err)
				} else {
					log.Fatal(err)
				}
			}
		}
	}
}

func deleteData(t Table, forcePurge bool, txDst *sql.Tx) error {
	log.Debug(fmt.Sprintf("DELETE data from %s according to configuration", t.Name))
	dst_query := fmt.Sprintf("DELETE FROM %s.%s", t.Schema, t.Name)
	_, err := txDst.Exec(dst_query)
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
func OrderTableList(table_list []string, txDst *sql.Tx) []string {

	var pkName string
	var nb_fk int
	var ordered_table_list []string
	tables := "{" + strings.Join(table_list, ",") + "}"

	// Query data from tableA
	rows, erri := txDst.Query(qry_linked_tables, tables)
	if erri != nil {
		log.Fatal("Error querying data from tableA: ", erri)
	}
	// Iterate through the rows from tableA and insert into tableB
	for rows.Next() {
		if erri := rows.Scan(&pkName, &nb_fk); erri != nil {
			log.Error("Error scanning row: ", erri)

		}
		ordered_table_list = append(ordered_table_list, pkName)
	}
	rows.Close()
	log.Debug(pkName)
	return ordered_table_list
}
