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

//go:embed sql/fetch_trigger_constraints.sql
var qry_fetch_trigger_constraints string

//go:embed sql/disable_trigger.sql
var qry_disable_trigger string

//go:embed sql/enable_trigger.sql
var qry_enable_trigger string

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

	rows, err := dbConn.Query(qry_fetch_trigger_constraints, filter)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}

	//	var tableFullname string
	//	var triggerName string
	//	var connName string
	var trgCons []TriggerConstraint
	var trx TriggerConstraint

	for rows.Next() {
		log.Debug(rows)
		rows.Scan(&trx.TableFullName, &trx.ConstraintName, &trx.TriggerName)
		log.Debug(fmt.Sprintf("%s", trx.TriggerName))

		trgCons = append(trgCons, trx)
		log.Debug(trx)
	}

	return trgCons
}

// Disable all triggers on database
func DisableTriggerConstraints(dbConn *sql.DB, triggers []TriggerConstraint) error {

	var err error
	for _, t := range triggers {

		err = DisableTrigger(dbConn, t.TableFullName, t.ConstraintName)
	}

	return err
}

// Disable a trigger on database
func DisableTrigger(dbConn *sql.DB, tableFullname string, triggerName string) error {
	log.Debug(fmt.Sprintf("%s DISABLE TRIGGER %s : ", tableFullname, triggerName))

	qry := fmt.Sprintf("ALTER TABLE %s ALTER CONSTRAINT %s DEFERRABLE", tableFullname, triggerName)

	// _, err := dbConn.Exec(qry_disable_trigger, tableFullname, triggerName)
	_, err := dbConn.Exec(qry)
	if err != nil {
		log.Debug(qry)
		log.Fatal("Error executing query:", err)
	}

	return err
}

// Enable all triggers on database
func EnableTriggerConstraints(dbConn *sql.DB, triggers []TriggerConstraint) error {

	var err error

	for _, t := range triggers {
		err = EnableTrigger(dbConn, t.TableFullName, t.TriggerName)
	}

	return err
}

// Enable a trigger on database
func EnableTrigger(dbConn *sql.DB, tableFullname string, triggerName string) error {
	log.Debug(fmt.Sprintf("%s ENABLE TRIGGER %s : ", tableFullname, triggerName))
	_, err := dbConn.Exec(qry_enable_trigger, tableFullname, triggerName)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}

	return err
}
