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

func insertMultiData(dbDst *sql.Tx, tableFullname string, colnames []string, colparam []string, multirows [][]interface{}) (int, string) {
	col_names := strings.Join(colnames, ",")

	nbRows := len(multirows)
	var errCode string
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

		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			log.Debug(pqErr.Code)
			errCode = fmt.Sprintf("%s", pqErr.Code)
		}

		log.Debug("Error :", err)
		log.Debug(destQuery)
		log.Warning(fmt.Sprintf("Error during INSERT on %s, retry", tableFullname))
		return 0, errCode
	}
	return 0, errCode
}
