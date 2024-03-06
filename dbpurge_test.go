package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableByNameExists(t *testing.T) {

	var tables []Table
	var columns []Column

	expectedValue := Table{"foobar", "public.foobar", columns, "public", "delete", "", ""}

	tables = append(tables, expectedValue)

	config := Config{tables, []Column{}}

	result, found := getTableByName(config, "foobar")

	assert.Equal(t, expectedValue, result, "Found table foobar")
	assert.Equal(t, found, true, "Found table foobar")
}

func TestGetTableByNameNotExists(t *testing.T) {

	var tables []Table
	var columns []Column

	expectedValue := Table{"foobar", "public.foobar", columns, "public", "delete", "", ""}

	tables = append(tables, expectedValue)

	config := Config{tables, []Column{}}

	result, found := getTableByName(config, "not_exists")

	assert.Equal(t, result, Table{}, "Table not found")
	assert.Equal(t, found, false, "Table not found")
}
