package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableByFullNameExists(t *testing.T) {
	table := Table{"foobar", "public.foobar", []Column{}, "public", "delete", "", ""}

	tables := []Table{table}

	config := Config{tables, []defColumn{}}

	result, found := getTableByFullName(config, "public.foobar")

	assert.Equal(t, result, table, "Found table foobar")
	assert.Equal(t, found, true, "Found table foobar")

}

func TestGetTableByFullNameNotExists(t *testing.T) {
	table := Table{"foobar", "public.foobar", []Column{}, "public", "delete", "", ""}

	tables := []Table{table}

	config := Config{tables, []defColumn{}}

	result, found := getTableByFullName(config, "public.barfoo")

	assert.Equal(t, result, Table{}, "Found table foobar")
	assert.Equal(t, found, false, "Found table foobar")

}

func TestGetTableByNameExists(t *testing.T) {

	var tables []Table
	var columns []Column

	expectedValue := Table{"foobar", "public.foobar", columns, "public", "delete", "", ""}

	tables = append(tables, expectedValue)

	config := Config{tables, []defColumn{}}

	result, found := getTableByName(config, "foobar")

	assert.Equal(t, expectedValue, result, "Found table foobar")
	assert.Equal(t, found, true, "Found table foobar")
}

func TestGetTableByNameNotExists(t *testing.T) {

	var tables []Table
	var columns []Column

	expectedValue := Table{"foobar", "public.foobar", columns, "public", "delete", "", ""}

	tables = append(tables, expectedValue)

	config := Config{tables, []defColumn{}}

	result, found := getTableByName(config, "not_exists")

	assert.Equal(t, result, Table{}, "Table not found")
	assert.Equal(t, found, false, "Table not found")
}

// Column defined to generate config file
// type dbForeignKey struct {
//	SchemaSource string
//	TableSource  string
//	SchemaTarget string
//	TableTarget  string
//	ColumnSource string
//	ColumnTarget string
// }

func TestQueryAddForeignKey(t *testing.T) {

	fk := dbForeignKey{"public", "book", "public", "author", "author_id", "id"}

	expected := "ALTER TABLE public.book ADD CONSTRAINT expulo_public_book_author_author_id_id_fkey FOREIGN KEY (author_id) REFERENCES public.author(id) ON DELETE CASCADE NOT VALID"

	result := queryAddForeignKey(fk)

	assert.Equal(t, result, expected)
}

func TestQueryDropForeignKey(t *testing.T) {

	fk := dbForeignKey{"public", "book", "public", "author", "author_id", "id"}

	expected := "ALTER TABLE public.book DROP CONSTRAINT expulo_public_book_author_author_id_id_fkey"

	result := queryDropForeignKey(fk)

	assert.Equal(t, result, expected)
}
