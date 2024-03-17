package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The column exist
func TestGetCols(t *testing.T) {

	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}

	conf := Table{"boat", "sea.boat", []Column{column}, "sea", "delete", "id < 42", ""}

	col, _ := getCols(conf, "id")

	if col != column {
		t.Fatalf("getCols does not return valid dsn")
	}

}

// The column does not exist
func TestGetColsNotFound(t *testing.T) {

	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 3, false}

	conf := Table{"boat", "sea.boat", []Column{column}, "sea", "delete", "id < 42", ""}

	_, found := getCols(conf, "name")

	if found {
		t.Fatalf("getCols does not return valid dsn")
	}

}

// The table exists
func TestCheckConfigTablesExists(t *testing.T) {

	configTables := []Table{{"boat", "sea.boat", []Column{}, "sea", "delete", "id < 42", ""}}
	existingTables := []string{"sea.skipper", "sea.boat"}

	existing, table := checkConfigTables(configTables, existingTables, "foobar")

	assert.Equal(t, true, existing, "The table exists")
	assert.Equal(t, "", table, "The table exists")

}

// The table does not exist
func TestCheckConfigTablesNotExists(t *testing.T) {

	configTables := []Table{{"boat", "sea.boat", []Column{}, "sea", "delete", "id < 42", ""}}
	existingTables := []string{"sea.fish"}

	existing, table := checkConfigTables(configTables, existingTables, "foobar")

	assert.Equal(t, false, existing, "The table does not exist")
	assert.Equal(t, "The table sea.boat does not exist in foobar database, check the configuration", table, "The table does not exist")

}

// The table does not exist
func TestCheckConfig(t *testing.T) {

	result := checkConfig(true, "foobar")

	assert.Equal(t, true, result, "The config is ok")
}

func TestGenerateConfig(t *testing.T) {
	file, _ := ioutil.TempFile("/tmp", "prefix")

	table := dbTable{"name", "schema", "clean_method", []dbColumn{{"name", "gen"}}}
	var tables []dbTable
	tables = append(tables, table)

	generateConfig(file.Name(), tables)

	jsonFile, _ := os.Open(file.Name())
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Config array
	var conf Config

	// we unmarshal our byteArray which contains our
	// jsonFile's content into our main Struct Config
	if err := json.Unmarshal(byteValue, &conf); err != nil {
		panic(err)
	}

	// assert.Equal(t, true, exists, "The config is ok")

	os.Remove(file.Name())
}

// func getDefaultGeneratorByName(name string) (bool, string)
//
func TestCheckConfigEmail(t *testing.T) {

	defaults := []defColumn{
		{"firstname", "FakeFirstName", false},
		{"name", "FakeName", false},
		{"email", "FakeEmail", false}}

	config.Defaults = defaults

	result, gen := getDefaultGeneratorByName("email")

	assert.Equal(t, true, result, "The config is ok")
	assert.Equal(t, gen, "FakeEmail", "The config is ok")
}

func TestCheckConfigFoobar(t *testing.T) {

	defaults := []defColumn{
		{"firstname", "FakeFirstName", false},
		{"name", "FakeName", false},
		{"email", "FakeEmail", false}}

	config.Defaults = defaults

	result, gen := getDefaultGeneratorByName("foobar")

	assert.Equal(t, result, false, "foobar does not exists")
	assert.Equal(t, gen, "", "foobar does not exists")
}

//
// func checkConfigGenerators(configTables []Table, allowed []string) (bool, string)
//
// all generators are allowed
func TestCheckConfigGeneratorsOK(t *testing.T) {
	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	tables := []Table{{"boat", "sea.boat", []Column{column}, "sea", "delete", "id < 42", ""}}

	generators := []string{"serial", "random"}

	result, gen := checkConfigGenerators(tables, generators)

	assert.Equal(t, result, true, "all generators are ok")
	assert.Equal(t, gen, "", "all generators are ok")
}

// all generator random does not exist
func TestCheckConfigGeneratorsKO(t *testing.T) {

	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	tables := []Table{{"boat", "sea.boat", []Column{column}, "sea", "delete", "id < 42", ""}}

	generators := []string{"serial"}

	result, gen := checkConfigGenerators(tables, generators)

	assert.Equal(t, result, false, "all generators are ok")
	assert.Equal(t, gen, "The generator random does not exist, check the configuration", "all generators are ok")
}
