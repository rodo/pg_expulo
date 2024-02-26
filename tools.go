package main

import (
	"fmt"
	"strings"
)

func toolPat(nbRows int, colparam []string) string {
	x := 1
	r := 0

	var i []string
	var j []string
	var c int
	var d int

	for r < nbRows {
		c = 0
		d = 1
		for c < len(colparam) {
			if colparam[c] != fmt.Sprintf("$%d", d) {
				j = append(j, colparam[c])
			} else {
				j = append(j, fmt.Sprintf("$%d", x))
				x++
				d++
			}
			c++
		}
		a := strings.Join(j, ",")
		i = append(i, fmt.Sprintf("(%s)", a))
		j = []string{}
		r++
	}

	return strings.Join(i, ",")
}
