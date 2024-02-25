package main

import (
	"testing"
)

func TestGet_cols(t *testing.T) {

	conf := Table { "boat" }

	_, err := get_cols(conf, "clean")

	if err != 0 {
		t.Fatalf("get_cols does not return valid dsn")
	}

}
