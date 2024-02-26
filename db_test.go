package main

import (
	"testing"
)

func TestGetDsn(t *testing.T) {

	_, dsn := getDsn("host", "port", "user", "pass", "db", "version")
	want := "user:pass@host:port/db"

	if dsn != want {
		t.Fatalf("getDsn does not retunr valid dsn")
	}
}

func TestFullTableName(t *testing.T) {

	result := fullTableName("schema", "name")
	want := "schema.name"

	if result != want {
		t.Fatalf("fullTableName does not return %s in place of %s", result, want)
	}
}
