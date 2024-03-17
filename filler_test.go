package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

// AllowedGenerators 2 cases
func TestAllowedGenerators(t *testing.T) {
	result := allowedGenerators()
	// Assert that result is of type []string
	assert.True(t, slices.Contains(result, "serial"), "Expected result of type string")
}

func TestAllowedGeneratorsFalse(t *testing.T) {
	result := allowedGenerators()
	// Assert that result is of type []string
	assert.False(t, slices.Contains(result, "foobar"), "Expected result of type string")
}

// setRandom 4 cases
func TestSetRandomTrue(t *testing.T) {
	result := setRandom(true, 10, 12)
	assert.Equal(t, result, 12, "Preserve null with non null value")
}

func TestSetRandomFalse(t *testing.T) {
	result := setRandom(false, 10, 12)
	assert.Equal(t, result, 12, "Do not preserve null with non null value")
}

func TestSetRandomNullTrue(t *testing.T) {
	result := setRandom(true, nil, 12)
	assert.Equal(t, result, nil, "Preserve null value with null value")
}

func TestSetRandomNullFalse(t *testing.T) {
	result := setRandom(false, nil, 17)
	assert.Equal(t, result, 17, "Do not preserve null with null value")
}

// fillColum
func TestFillColumnIgnore(t *testing.T) {
	var colValues []interface{}
	var colparam []string
	var cols []interface{}
	var colNames []string
	var columns []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}

	i := 0
	nbColumnModified := 0
	cfvalue := "ignore"

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 0, "No column modified")
}

// null
func TestFillColumnNull(t *testing.T) {
	var colValues []interface{}
	var colparam []string
	var cols []interface{}
	var colNames []string

	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}

	i := 0

	nbColumnModified := 1
	cfvalue := "null"
	columns := []string{"foo"}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 2, "One column is modified")
}

// mask
func TestFillColumnMask(t *testing.T) {
	var colValues []interface{}
	var colparam []string
	var colNames []string

	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}

	i := 0

	nbColumnModified := 1
	cfvalue := "mask"
	columns := []string{"foo"}
	cols := []interface{}{"foo", 3}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 2, "One column is modified")
	assert.Equal(t, colValues, []interface{}{"********"}, "One column is modified")
}

// md5
func TestFillColumnMd5(t *testing.T) {
	var colValues []interface{}
	var colparam []string
	var colNames []string

	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}

	i := 0

	nbColumnModified := 1
	cfvalue := "md5"
	columns := []string{"foo"}
	cols := []interface{}{"foo", 3}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 2, "One column is modified")
	assert.Equal(t, colValues, []interface{}{"acbd18db4cc2f85cedef654fccc4a4d8"}, "One column is modified")
}

// sql
func TestFillColumnSql(t *testing.T) {
	var colValues []interface{}
	var colparam []string
	var colNames []string

	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}

	i := 0

	nbColumnModified := 1
	cfvalue := "sql"
	columns := []string{"foo"}
	cols := []interface{}{"foo", 3}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 1, "One column is modified")
	assert.Equal(t, colparam, []string{"getRandomString()"}, "One column is modified")
}
