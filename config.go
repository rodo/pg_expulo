package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// Config store the whole configuration read from json file
type Config struct {
	Tables []Table `json:"tables"`
}

// Columns contains a collection of Column
type Columns struct {
	Columns []Column `json:"columns"`
}

// Table represent a table with her property in configuration file
type Table struct {
	Name        string   `json:"name"`
	Columns     []Column `json:"columns"`
	Schema      string   `json:"schema"`
	CleanMethod string   `json:"clean"`
	Filter      string   `json:"filter"`
}

type Column struct {
	Name        string `json:"name"`
	Generator   string `json:"generator"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Timezone    string `json:"timezone"`
	SQLFunction string `json:"function"`
}

func read_config(fileName string) Config {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	log.Info("Successfully Opened : ", fileName)
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

	// tlist := []string{}

	// for i := 0; i < len(tables.Tables); i++ {
	//	tlist = append(tlist, tables.Tables[i].Name)

	//	for j := 0; j < len(tables.Tables[i].Columns); j++ {

	//		fullname := fmt.Sprintf("%s.%s.%s", tables.Tables[i].Schema,
	//			tables.Tables[i].Name,
	//			tables.Tables[i].Columns[j].Name)
	//		//log.Debug(fullname)
	//	}
	// }

	return tables
}

// Return the column columnName in the table Table
func get_cols(conf Table, columName string) (Column, bool) {
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
