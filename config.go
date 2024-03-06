package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

// Read configuration file in json from disk
func readConfig(filename string) Config {

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("file '%s' does not exist", filename))
	}

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("error opening file '%s': %v", filename, err))
	}
	log.Info("Successfully Opened : ", filename)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	// we initialize our Config array
	var conf Config

	// we unmarshal our byteArray which contains our
	// jsonFile's content into our main Struct Config
	if err := json.Unmarshal(byteValue, &conf); err != nil {
		panic(err)
	}

	return conf
}

// Extend the confirmation with sequence information
// Set the information on each column when they are defined as serial
func getInfoFromDatabases(config Config, sequences []Sequence) Config {

	var tables []Table

	for _, t := range config.Tables {
		var newColumns []Column
		t.FullName = fmt.Sprintf("%s.%s", t.Schema, t.Name)
		for _, c := range t.Columns {
			newColumn := c

			for _, v := range sequences {
				if t.FullName == v.TableName && c.Name == v.ColumnName {
					log.Debug(fmt.Sprintf("Assign seq last value %d to %s.%s based on %s", v.LastValue, v.TableName, newColumn.Name, v.SequenceName))
					newColumn.SequenceName = v.SequenceName
					newColumn.SeqLastValue = int64(v.LastValue)
				}
			}
			//
			newColumns = append(newColumns, newColumn)
			t.Columns = newColumns
		}
		tables = append(tables, t)
	}

	newconf := Config{tables}

	return newconf
}

// Return the column columnName in the table Table
func getCols(conf Table, columName string) (Column, bool) {
	found := false
	var result Column
	for j := 0; j < len(conf.Columns); j++ {
		if columName == conf.Columns[j].Name {
			result = conf.Columns[j]
			found = true
		}

	}
	return result, found
}

// Check if all tables defined in configuration exist in the database
func checkConfigTables(configTables []Table, existingTables []string, dbname string) (bool, string) {
	result := true
	var tf string
	for _, t := range configTables {
		tf = fullTableName(t.Schema, t.Name)
		if slices.Contains(existingTables, tf) {
			result = true
		} else {
			return false, fmt.Sprintf("The table %s does not exist in %s database, check the configuration", tf, dbname)
		}
	}
	return result, ""
}

// Check if generators are well written
func checkConfigGenerators(configTables []Table, allowed []string) (bool, string) {
	result := true

	for _, t := range configTables {
		for _, c := range t.Columns {
			if slices.Contains(allowed, c.Generator) {
				result = true
			} else {
				return false, fmt.Sprintf("The generator %s does not exist, check the configuration", c.Generator)
			}
		}
	}
	return result, ""
}

func checkConfig(result bool, message string) bool {

	if !result {
		log.Fatal(message)
	}
	return result
}
