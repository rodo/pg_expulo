package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
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

	// we initialize our Tables array
	var tables Config

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'tables' which we defined above
	if err := json.Unmarshal(byteValue, &tables); err != nil {
		panic(err)
	}

	return tables
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
