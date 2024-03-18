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

// a Fake function
func TestFillColumnFakeEmail(t *testing.T) {

	cfvalue := "FakeEmail"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 2
	columns := []string{"foo"}
	cols := []interface{}{"foo", 0}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 3, "One column is modified")
	assert.Equal(t, colparam, []string{"$2"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.IsType(t, colValues[0], "astring", "One column is modified")
}

// there is no generator defined for this column
// there is no default defined
func TestFillColumnNotFoundNoDefault(t *testing.T) {

	cfvalue := "notfound"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"name", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 2
	columns := []string{"foobar"}
	cols := []interface{}{"foobar", 0}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 3, "One column is modified")
	assert.Equal(t, colparam, []string{"$2"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.Equal(t, colValues[0], "foobar", "the value is set to foobar")
	assert.IsType(t, colValues[0], "astring", "the result is a string")
}

// there is no generator defined for this column
// there is a default value in Defaults for this column
func TestFillColumnNotFoundDefault(t *testing.T) {

	cfvalue := "notfound"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"name", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 2
	columns := []string{"name"}
	cols := []interface{}{"foobar", 0}

	// We define default values for the columns named "name"
	defaults := []defColumn{
		{"name", "FakeName", false},
		{"email", "FakeEmail", false}}

	config.Defaults = defaults

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 3, "One column is modified")
	assert.Equal(t, colparam, []string{"$2"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.NotEqual(t, colValues[0], "foobar", "One column is modified")
	assert.IsType(t, colValues[0], "astring", "One column is modified")
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

// randomInt
func TestFillColumnRandomInt(t *testing.T) {
	// the case we want to test
	cfvalue := "randomInt"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 10
	columns := []string{"foo"}
	cols := []interface{}{"foo", 0}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 11, "One column is modified")
	assert.Equal(t, colparam, []string{"$10"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.IsType(t, colValues[0], int32(10.0), "One column is modified")
}

// randomInt32
func TestFillColumnRandomInt32(t *testing.T) {
	// the case we want to test
	cfvalue := "randomInt32"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 10
	columns := []string{"foo"}
	cols := []interface{}{"foo", 0}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 11, "One column is modified")
	assert.Equal(t, colparam, []string{"$10"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.IsType(t, colValues[0], int32(10.0), "One column is modified")
}

// randomFloat
func TestFillColumnRandomFloat(t *testing.T) {
	// the case we want to test
	cfvalue := "randomFloat"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 10
	columns := []string{"foo"}
	cols := []interface{}{"foo", 0}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 11, "One column is modified")
	assert.Equal(t, colparam, []string{"$10"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.IsType(t, colValues[0], float32(10.0), "One column is modified")
}

// randomFloat64
func TestFillColumnRandomFloat64(t *testing.T) {
	// the case we want to test
	cfvalue := "randomFloat64"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 10
	columns := []string{"foo"}
	cols := []interface{}{"foo", 0}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 11, "One column is modified")
	assert.Equal(t, colparam, []string{"$10"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.IsType(t, colValues[0], float64(10.0), "One column is modified")
}

// randomString
func TestFillColumnRandomString(t *testing.T) {

	cfvalue := "randomString"

	var colValues []interface{}
	var colparam []string
	var colNames []string
	var sequences map[string]Sequence
	var foreignKeys map[string]string
	var initValues map[string]int64

	c := Column{"id", "random", 0, 42, "UTC", "getRandomString()", "id_seq", 1, false}
	a := Table{"boat", "sea.boat", []Column{c}, "sea", "delete", "id < 42", ""}
	i := 0

	nbColumnModified := 2
	columns := []string{"foo"}
	cols := []interface{}{"foo", 0}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 3, "One column is modified")
	assert.Equal(t, colparam, []string{"$2"}, "One column is modified")
	assert.Equal(t, len(colValues), 1, "One column is modified")
	assert.IsType(t, colValues[0], "astring", "One column is modified")
	assert.NotEqual(t, colValues[0], "foo", "The content of the row is modified")
}

// md5
func TestFillColumnMd5(t *testing.T) {

	cfvalue := "md5"

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
	columns := []string{"foo"}
	cols := []interface{}{"foo", 3}

	fillColumn(a, c, cfvalue, &colValues, &colparam, &nbColumnModified, cols, &colNames, i, columns, &sequences, foreignKeys, initValues)

	assert.Equal(t, nbColumnModified, 2, "One column is modified")
	assert.Equal(t, colValues, []interface{}{"acbd18db4cc2f85cedef654fccc4a4d8"}, "One column is modified")
	assert.NotEqual(t, colValues[0], "foo", "The content of the row is modified")
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
