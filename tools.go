package main

import (
	"fmt"
	"strings"
)

func toolPat(nbRows int, nbCols int) string {
	x := 1
	r := 0

	var i []string
	var j []string
	c := 0
	for r < nbRows {
		for c < nbCols {

			j = append(j, fmt.Sprintf("$%d", x))
			c = c + 1
			x = x + 1

		}
		a := strings.Join(j, ",")
		i = append(i, fmt.Sprintf("(%s)", a))

		c = 0
		j = []string{}
		r = r + 1
	}

	return strings.Join(i, ",")
}
