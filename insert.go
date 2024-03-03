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

// Build the statement for bulk insert
func prepareQuery(tableName string, colNames []string, pat string) string {
	cols := strings.Join(colNames, ",")
	qrs := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, cols, pat)
	return qrs
}

func prepareValues(multirows [][]interface{}) []interface{} {
	var allValues []interface{}
	for _, row := range multirows {
		// Append each element of the row to allValues
		allValues = append(allValues, row...)
	}
	return allValues
}

func insertMultiData(dbDst *sql.Tx, tableFullname string, colNames []string, stmtParam []string, multirows [][]interface{}) (int, string) {
	var errCode string

	pat := toolPat(len(multirows), stmtParam)
	destQuery := prepareQuery(tableFullname, colNames, pat)
	allValues := prepareValues(multirows)

	_, err := dbDst.Exec(destQuery, allValues...)
	if err != nil {

		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			log.Debug(pqErr.Code)
			errCode = string(pqErr.Code)
		}

		log.Debug("Error :", err)
		log.Debug(destQuery)
		log.Warning(fmt.Sprintf("Error during INSERT on %s, retry", tableFullname))
		return 0, errCode
	}
	return 0, errCode
}
