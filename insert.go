package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func insertMultiData(dbDst *sql.Tx, tableFullname string, colnames []string, colparam []string, multirows [][]interface{}) {
	col_names := strings.Join(colnames, ",")

	nbRows := len(multirows)

	pat := toolPat(nbRows, colparam)

	// log.Debug(fmt.Sprintf("there is %d rows of %d columns", nbRows, nbColumns))

	var allValues []interface{}
	for _, row := range multirows {
		// Append each element of the row to allValues
		allValues = append(allValues, row...)
	}

	destQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableFullname, col_names, pat)

	_, err := dbDst.Exec(destQuery, allValues...)
	if err != nil {
		log.Debug("Error :", err)
		log.Debug(destQuery)
		log.Debug(allValues)
		log.Fatal("Error during INSERT on : ", tableFullname)

		return
	}
}
