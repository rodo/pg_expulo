package main

import (
	"fmt"
	"reflect"

	"github.com/go-faker/faker/v4"
)

type R struct{}

func (R) FakeEmail() string     { return faker.Email() }
func (R) FakeName() string      { return faker.Name() }
func (R) FakeFirstName() string { return faker.FirstName() }

func FillColumn(col Column, cfvalue string, colvalue []interface{}, colparam []string, nbcol int, cols []interface{}, colnames []string, i int, columns []string, lastvalue int64) ([]interface{}, []string, int, []string, int64) {

	x := int64(0)

	// The column is ignored in configuration
	if cfvalue == "ignore" {
		return colvalue, colparam, nbcol, colnames, x
	}

	// The column is NOT ignored in configuration
	colnames = append(colnames, columns[i])

	longv := fmt.Sprintf("%s******", cfvalue)
	if longv[:4] == "Fake" {
		v := reflect.ValueOf(R{})
		m := v.MethodByName(cfvalue)
		res := m.Call(nil)

		colvalue = append(colvalue, fmt.Sprintf("%s", res[0]))
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))

		nbcol++

		return colvalue, colparam, nbcol, colnames, x
	}

	// Assign the target value
	switch cfvalue {
	case "serial":
		// Set the column with NULL values
		if val, ok := cols[i].(int64); ok {
			// If it is, perform the addition
			x = val + lastvalue
		}
		colvalue = append(colvalue, x)
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))

	case "null":
		// Set the column with NULL values
		colvalue = append(colvalue, nil)
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "mask":
		colvalue = append(colvalue, mask())
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "randomInt":
		colvalue = append(colvalue, randomInt())
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "randomIntMinMax":
		colvalue = append(colvalue, randomIntMinMax(col.Min, col.Max))
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "randomFloat":
		colvalue = append(colvalue, randomFloat())
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "randomString":
		colvalue = append(colvalue, randomString())
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "md5":
		colvalue = append(colvalue, md5signature(fmt.Sprintf("%v", cols[i])))
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "randomTimeTZ":
		colvalue = append(colvalue, randomTimeTZ(col.Timezone))
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	case "sql":
		nbcol--
		colparam = append(colparam, col.SQLFunction)
	default:
		colvalue = append(colvalue, cols[i])
		colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
	}
	nbcol++

	return colvalue, colparam, nbcol, colnames, x
}
