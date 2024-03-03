package main

import (
	"fmt"
	"reflect"

	"github.com/go-faker/faker/v4"
)

type genericFake struct{}

func (genericFake) FakeEmail() string     { return faker.Email() }
func (genericFake) FakeName() string      { return faker.Name() }
func (genericFake) FakeFirstName() string { return faker.FirstName() }

//gocyclo:ignore
func fillColumn(table Table, col Column, cfvalue string, colValues *[]interface{}, colparam *[]string, nbColumnModified *int, cols []interface{}, colNames *[]string, i int, columns []string, sequences *map[string]Sequence, foreignKeys map[string]string, initValues map[string]int64) {

	x := int64(0)

	// The column is ignored in configuration
	if cfvalue == "ignore" {
		return
	}

	// The column is NOT ignored in configuration
	*colNames = append(*colNames, columns[i])

	longv := fmt.Sprintf("%s******", cfvalue)
	if longv[:4] == "Fake" {
		v := reflect.ValueOf(genericFake{})
		m := v.MethodByName(cfvalue)
		res := m.Call(nil)

		*colValues = append(*colValues, fmt.Sprintf("%s", res[0]))
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))

		*nbColumnModified++

		return
	}

	// Assign the target value
	switch cfvalue {
	case "foreign_key":

		colkey := fmt.Sprintf("%s.%s", table.FullName, col.Name)
		valkey := foreignKeys[colkey]

		val := initValues[valkey]

		// Deal with null in foreign key
		if _, ok := cols[i].(int64); ok {
			*colValues = append(*colValues, cols[i].(int64)+val)
		} else {
			*colValues = append(*colValues, cols[i])
		}

		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "serial":
		// Set the column with NULL values
		if val, ok := cols[i].(int64); ok {
			// If it is, perform the addition
			x = val + col.SeqLastValue

			seq := (*sequences)[col.SequenceName]
			seq.LastValueUsed = x

			(*sequences)[col.SequenceName] = seq
		}

		*colValues = append(*colValues, x)
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))

	case "null":
		// Set the column with NULL values
		*colValues = append(*colValues, nil)
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "mask":
		*colValues = append(*colValues, mask())
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "randomInt":
		*colValues = append(*colValues, randomInt())
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "randomInt32":
		*colValues = append(*colValues, randomInt32())
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "randomFloat64":
		*colValues = append(*colValues, randomFloat64())
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "randomIntMinMax":
		*colValues = append(*colValues, randomIntMinMax(col.Min, col.Max))
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "randomFloat":
		*colValues = append(*colValues, randomFloat())
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "randomString":
		*colValues = append(*colValues, randomString())
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "md5":
		*colValues = append(*colValues, md5signature(fmt.Sprintf("%v", cols[i])))
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "randomTimeTZ":
		*colValues = append(*colValues, randomTimeTZ(col.Timezone))
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	case "sql":
		*nbColumnModified--
		*colparam = append(*colparam, col.SQLFunction)
	default:
		*colValues = append(*colValues, cols[i])
		*colparam = append(*colparam, fmt.Sprintf("$%d", *nbColumnModified))
	}
	*nbColumnModified++

}
