package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The column exist
func TestGetCols(t *testing.T) {

	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1}

	conf := Table{"boat", "sea.boat", []Column{column}, "sea", "delete", "id < 42", ""}

	col, _ := getCols(conf, "id")

	if col != column {
		t.Fatalf("getCols does not return valid dsn")
	}

}

// The column does not exist
func TestGetColsNotFound(t *testing.T) {

	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 3}

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

	existing, table := checkConfigTables(configTables, existingTables)

	assert.Equal(t, true, existing, "The table exists")
	assert.Equal(t, "", table, "The table exists")

}

// The table does not exist
func TestCheckConfigTablesNotExists(t *testing.T) {

	configTables := []Table{{"boat", "sea.boat", []Column{}, "sea", "delete", "id < 42", ""}}
	existingTables := []string{"sea.fish"}

	existing, table := checkConfigTables(configTables, existingTables)

	assert.Equal(t, false, existing, "The table does not exist")
	assert.Equal(t, "sea.boat", table, "The table does not exist")

}
