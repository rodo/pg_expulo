package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
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

func flagParse() {
	// Global variables
	flag.BoolVar(&tryOnly, "try", false, "ROLLBACK everything on target. No data will be inserted")

	flag.BoolVar(&purgeOnly, "purge", false, "Only purge destination tables and exit, no data will be inserted")

	flag.StringVar(&configFile, "config", "config.json", "Configuration file to use")

	// Local usage only
	debug := flag.Bool("debug", false, "run in loglevel DEBUG")
	version := flag.Bool("version", false, "display version")

	// Parse flags
	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}
