package main

import (
	"testing"
)

func TestToolPat(t *testing.T) {

	result := toolPat(2, 2, []string{"$1", "$2"})
	want := "($1,$2),($3,$4)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

func TestToolPatSix(t *testing.T) {

	result := toolPat(1, 3, []string{"$1", "$2", "$3"})
	want := "($1,$2,$3)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

func TestToolPatNine(t *testing.T) {

	result := toolPat(3, 3, []string{"$1", "$2", "$3"})
	want := "($1,$2,$3),($4,$5,$6),($7,$8,$9)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

func TestToolPatFunc(t *testing.T) {

	result := toolPat(2, 3, []string{"$1", "now()", "$3"})
	want := "($1,now(),$3),($4,now(),$6)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}
