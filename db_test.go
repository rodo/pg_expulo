package main

import (
	"testing"
)

func TestGet_dsn(t *testing.T) {

	_, dsn := get_dsn("host", "port", "user", "pass", "db", "version")
	want := "user:pass@host:port/db"

	if dsn != want {
		t.Fatalf("get_dsn does not retunr valid dsn")
	}
}

func TestFullTableName(t *testing.T) {

	result := fullTableName("schema", "name")
	want := "schema.name"

	if result != want {
		t.Fatalf("fullTableName does not return %s in cas of %s", result, want)
	}
}
