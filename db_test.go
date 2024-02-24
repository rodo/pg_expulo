package main

import (
	"testing"
)



func TestGet_dsn(t *testing.T) {

	_, dsn := get_dsn("host", "port", "user", "pass", "db", "version")

	want := "user:pass@host:port/db"


	if dsn != want {
		t.Fatalf("randomFloat32() does not return float32")
	}

}
