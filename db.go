package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

//go:embed sql/fetch_table_foreign_keys.sql
var qry_fetch_table_foreign_keys string

// Restart a value
func ResetSeq(dbConn *sql.Tx, seq string, newvalue int64) {
	query := "ALTER SEQUENCE %s RESTART WITH %d"

	qry := fmt.Sprintf(query, seq, newvalue+1)

	log.Debug(qry)

	_, err := dbConn.Exec(qry)
	if err != nil {
		log.Fatal("Error executing query in GetSeqLastValue:", err)
	}
}

// Retreive the last used value in a sequence
func GetSeqLastValue(dbConn *sql.Tx, seq string) (int64, error) {
	var err error

	query := "SELECT last_value FROM %s"

	var last_value int64
	last_value = 0

	qry := fmt.Sprintf(query, seq)

	log.Debug(qry)

	rows, err := dbConn.Query(qry)
	if err != nil {
		log.Fatal("Error executing query in GetSeqLastValue:", err)
	}

	for rows.Next() {
		err = rows.Scan(&last_value)
		if err != nil {
			log.Fatal("Error on row", err)
		}
		log.Debug(fmt.Sprintf("row values : %s %d", seq, last_value))

	}
	if err = rows.Err(); err != nil {
		log.Fatal("Error reading rows :", err)
	}
	rows.Close()

	return last_value, err
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

func fullname(schemaname string, datname string, attname string) string {
	return fmt.Sprintf("%s.%s.%s", schemaname, datname, attname)
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
func GetTriggerConstraints(dbConn *sql.DB, tbFullnames []string) []TriggerConstraint {

	filter := "{" + strings.Join(tbFullnames, ",") + "}"

	log.Debug("filter : ", filter)

	rows, err := dbConn.Query(qry_fetch_table_foreign_keys, filter)
	if err != nil {
		log.Fatal("Error executing query in GetTriggerConstraints:", err)
	}

	var trgCons []TriggerConstraint
	var trx TriggerConstraint

	for rows.Next() {
		err = rows.Scan(&trx.TableFullName, &trx.ConstraintName)
		if err != nil {
			log.Fatal("Error on row", err)
		}
		log.Debug(fmt.Sprintf("row values : %s %s", trx.TableFullName, trx.ConstraintName))

		trgCons = append(trgCons, trx)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("Error reading rows :", err)
	}
	rows.Close()

	return trgCons
}

// Disable all triggers on database
func DeferForeignKeys(dbConn *sql.DB, triggers []TriggerConstraint) error {

	var err error

	qry := "ALTER TABLE %s ALTER CONSTRAINT %s INITIALLY DEFERRED"

	for _, t := range triggers {
		err = AlterForeignKey(dbConn, qry, t.TableFullName, t.ConstraintName)
	}

	return err
}

// Reactivate all foreign keys
func ReactivateForeignKeys(dbConn *sql.DB, triggers []TriggerConstraint) error {

	var err error

	qry := "ALTER TABLE %s ALTER CONSTRAINT %s NOT DEFERRABLE"

	for _, t := range triggers {
		err = AlterForeignKey(dbConn, qry, t.TableFullName, t.ConstraintName)
	}

	return err
}

// Disable a trigger on database
func AlterForeignKey(dbConn *sql.DB, queryDef string, tablename string, fkName string) error {
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
func CloseTx(tx *sql.Tx, tryOnly bool) string {

	if tryOnly {
		// Rollback the transaction on target as requested
		if err := tx.Rollback(); err != nil {
			log.Fatal("Error committing transaction: ", err)
		}
		log.Info("Rollback on target")
		return "rollback"
	} else {
		// Commit the transaction on target if all queries succeed
		if err := tx.Commit(); err != nil {
			log.Fatal("Error committing transaction: ", err)
		} else {
			log.Info("Commit on target")
		}
		return "commit"
	}
}
