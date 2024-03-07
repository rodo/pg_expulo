package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

//go:embed sql/fetch_foreign_keys.sql
var qryFetchForeignKeys string

//go:embed sql/tables.sql
var qryTables string

//go:embed sql/fetch_tables.sql
var qryFetchTables string

// Restart all the sequences
func resetAllSequences(dbConn *sql.DB, sequences *map[string]Sequence) {
	for _, s := range *sequences {
		if s.LastValueUsed > s.InitialValue {
			resetSeq(dbConn, s.SequenceName, s.LastValueUsed)
		}
	}
}

// Return an array of table name in fullname
func getExistingTables(dbConn *sql.DB) []string {
	var tables []string

	rows, err := dbConn.Query(qryTables)
	if err != nil {
		log.Fatal("Error executing query in GetExistingTables:", err)
	}

	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		tables = append(tables, table)
	}

	return tables
}

// Return an array of table name in fullname
func getDbTables(dbConn *sql.DB) []dbTable {
	var tables []dbTable

	rows, err := dbConn.Query(qryFetchTables)
	if err != nil {
		log.Fatal("Error executing query in GetExistingTables:", err)
	}

	for rows.Next() {
		var table dbTable
		err = rows.Scan(&table.Schema, &table.Name)
		table.CleanMethod = "delete"
		tables = append(tables, table)
	}

	return tables
}

// Restart a sequence with a new value
func resetSeq(dbConn *sql.DB, seq string, newvalue int64) {
	query := "ALTER SEQUENCE %s RESTART WITH %d"

	qry := fmt.Sprintf(query, seq, newvalue)

	log.Info(fmt.Sprintf("Restart sequence %s with value %d", seq, newvalue))

	_, err := dbConn.Exec(qry)
	if err != nil {
		log.Fatal("Error executing query in GetSeqLastValue:", err)
	}
}

func getDsn(host string, port string, user string, pass string, db string, version string) (string, string) {

	appname := "expulo_" + version

	// The connection string is used by Open() method
	cnx := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s application_name=%s sslmode=disable", host, port, user, pass, db, appname)

	// The dsn is used in log, as it's more readable
	dsn := fmt.Sprintf("%s:%s@%s:%s/%s", user, pass, host, port, db)
	return cnx, dsn
}

func connectDb(connectionString string) *sql.DB {
	// Connect to the database source
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Error connecting to the database source:", err)
		os.Exit(1)
	}

	// Ensure the connection is up and running
	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func fullTableName(schemaname string, datname string) string {
	return fmt.Sprintf("%s.%s", schemaname, datname)
}

func queryTableSource(dbSrc *sql.DB, query string) (*sql.Rows, []string) {

	rows, err := dbSrc.Query(query)
	log.Debug("Executing query : ", query)
	if err != nil {
		log.Fatal("Error executing query:", err)
		os.Exit(1)
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("Error getting column names:", err)
	}

	return rows, columns
}

// Fetch the trigger constraint from the database catalog
func getTriggerConstraints(dbConn *sql.DB, tbFullnames []string, foreignKeys *map[string]string) []TriggerConstraint {

	filter := "{" + strings.Join(tbFullnames, ",") + "}"

	log.Debug("filter : ", filter)

	rows, err := dbConn.Query(qryFetchForeignKeys, filter)
	if err != nil {
		log.Fatal("Error executing query in GetTriggerConstraints:", err)
	}

	var trgCons []TriggerConstraint
	var trx TriggerConstraint
	var tableTarget string
	// The id of the column on pg_attribute
	var columnTarget string
	var columnSource string

	for rows.Next() {
		err = rows.Scan(&trx.TableFullName, &tableTarget, &trx.ConstraintName, &columnTarget, &columnSource)
		if err != nil {
			log.Fatal("Error on row", err)
		}
		log.Debug(fmt.Sprintf("row values : %s %s %s %s", trx.TableFullName, trx.ConstraintName, columnTarget, columnSource))

		(*foreignKeys)[columnTarget] = columnSource

		trgCons = append(trgCons, trx)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("Error reading rows :", err)
	}
	rows.Close()

	return trgCons
}

// Disable all triggers on database
func deferForeignKeys(dbConn *sql.DB, triggers []TriggerConstraint) error {

	var err error

	qry := "ALTER TABLE %s ALTER CONSTRAINT %s INITIALLY DEFERRED"

	for _, t := range triggers {
		err = alterForeignKey(dbConn, qry, t.TableFullName, t.ConstraintName)
	}

	return err
}

// Reactivate all foreign keys
func reactivateForeignKeys(dbConn *sql.DB, triggers []TriggerConstraint) error {

	var err error

	qry := "ALTER TABLE %s ALTER CONSTRAINT %s NOT DEFERRABLE"

	for _, t := range triggers {
		err = alterForeignKey(dbConn, qry, t.TableFullName, t.ConstraintName)
	}

	return err
}

// Disable a trigger on database
func alterForeignKey(dbConn *sql.DB, queryDef string, tablename string, fkName string) error {
	var err error

	qry := fmt.Sprintf(queryDef, tablename, fkName)

	log.Debug(qry)

	_, err = dbConn.Exec(qry)
	if err != nil {
		log.Fatal("Error executing query in DeferForeignKey:", err)
	}

	return err
}

// Commit or Roll back an open transaction
func closeTx(tx *sql.Tx, tryOnly bool) string {

	if tryOnly {
		// Rollback the transaction on target as requested
		if err := tx.Rollback(); err != nil {
			log.Fatal("Error committing transaction: ", err)
		}
		log.Info("Rollback on target")
		return "rollback"
	}

	// Commit the transaction on target if all queries succeed
	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction : ", err)
	} else {
		log.Info("Commit on target")
	}
	return "commit"
}

//go:embed sql/fetch_sequences.sql
var qryFetchSequences string

// Read the informations sequences from database
// Retreive the last used value in a sequence
func getSequencesInfo(dbConn *sql.DB) ([]Sequence, map[string]Sequence) {
	var err error
	var sequences []Sequence

	rows, err := dbConn.Query(qryFetchSequences)
	if err != nil {
		log.Fatal("Error executing query in GetSequencesInfo:", err)
	}

	for rows.Next() {
		var s = Sequence{}
		err = rows.Scan(&s.TableName, &s.ColumnName, &s.SequenceName, &s.LastValue, &s.ColumnPosition)
		s.InitialValue = int64(s.LastValue)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error on row in GetSequencesInfo %s", err))
		}
		sequences = append(sequences, s)
		log.Debug(fmt.Sprintf("Found sequence %s defined in database", s.SequenceName))
	}
	if err = rows.Err(); err != nil {
		log.Fatal("Error reading rows :", err)
	}
	rows.Close()

	log.Debug(fmt.Sprintf("Found %d sequence(s) defined in database", len(sequences)))

	mapSeq := make(map[string]Sequence)

	for _, seqX := range sequences {
		mapSeq[seqX.SequenceName] = seqX
	}

	// TODO refacto this part and return only a map instead of two
	return sequences, mapSeq
}
