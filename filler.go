package main

import "fmt"

func fillColumn(col Column, cfvalue string, colvalue []interface{}, colparam []string, nbcol int, cols []interface{}, colnames []string, i int, columns []string) ([]interface{}, []string, int, []string) {

	if cfvalue != "ignore" {

		colnames = append(colnames, columns[i])

		// Assign the target value
		switch cfvalue {
		case "null":
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
		case "fake_email":
			colvalue = append(colvalue, fakeEmail())
			colparam = append(colparam, fmt.Sprintf("$%d", nbcol))
		case "fake_name":
			colvalue = append(colvalue, fakeName())
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
	}

	return colvalue, colparam, nbcol, colnames

}
