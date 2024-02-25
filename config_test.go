package main

import (
	"testing"
)

// The column exist
func TestGet_cols(t *testing.T) {

	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()"}

	conf := Table{"boat", []Column{column}, "public", "delete", "id < 42"}

	col, _ := get_cols(conf, "id")

	if col != column {
		t.Fatalf("get_cols does not return valid dsn")
	}

}

// The column does not exist
func TestGet_colsNotFound(t *testing.T) {

	column := Column{"id", "random", 0, 42, "UTC", "getRandomString()"}

	conf := Table{"boat", []Column{column}, "public", "delete", "id < 42"}

	_, found := get_cols(conf, "name")

	if found {
		t.Fatalf("get_cols does not return valid dsn")
	}

}
