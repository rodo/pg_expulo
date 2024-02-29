package main

import (
	"testing"
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
