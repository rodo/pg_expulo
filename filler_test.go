package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestFakeEmail(t *testing.T) {
	result := genericFake{}.FakeEmail()
	assert.NotNil(t, result, "Expected non-nil result")

	// Assert that result is of type string
	assert.IsType(t, "", result, "Expected result of type string")
}

func TestFakeName(t *testing.T) {
	result := genericFake{}.FakeName()
	assert.NotNil(t, result, "Expected non-nil result")

	// Assert that result is of type string
	assert.IsType(t, "", result, "Expected result of type string")
}

func TestFakeFirstName(t *testing.T) {
	result := genericFake{}.FakeFirstName()

	assert.NotNil(t, result, "Expected non-nil result")

	// Assert that result is of type string
	assert.IsType(t, "", result, "Expected result of type string")
}

func TestAllowedGenerators(t *testing.T) {
	result := allowedGenerators()
	// Assert that result is of type []string
	assert.True(t, slices.Contains(result, "serial"), "Expected result of type string")
}

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
