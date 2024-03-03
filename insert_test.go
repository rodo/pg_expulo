package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareQuery(t *testing.T) {

	expectedValue := "INSERT INTO boat (id,names) VALUES ?,?"

	result := prepareQuery("boat", []string{"id", "names"}, "?,?")

	assert.Equal(t, expectedValue, result, "Found table foobar")

}

func TestPrepareValues(t *testing.T) {

	expectedValue := []interface{}{1, "2", 3, 4, "5", 6}

	result := prepareValues([][]interface{}{{1, "2", 3}, {4, "5", 6}})

	assert.Equal(t, expectedValue, result, "Wrong values")

}
