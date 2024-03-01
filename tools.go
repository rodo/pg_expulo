package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ConnectionParameters define all the parameters needed to connect to PostgreSQL
type ConnectionParameters struct {
	Host string
	Port string
	User string
	Pass string
	Db   string
}

func getEnvVariable(envName string) string {

	value, exists := os.LookupEnv(envName)

	if !exists {
		log.Error(fmt.Sprintf("Environment variable %s is not defined, please define it\n", envName))
		log.Fatal("exit on previous error")
	}

	return value
}

func readEnv(target string) ConnectionParameters {

	srcHost := getEnvVariable(fmt.Sprintf("PG%sHOST", target))
	srcPort := getEnvVariable(fmt.Sprintf("PG%sPORT", target))
	srcUser := getEnvVariable(fmt.Sprintf("PG%sUSER", target))
	srcPass := getEnvVariable(fmt.Sprintf("PG%sPASSWORD", target))
	srcDb := getEnvVariable(fmt.Sprintf("PG%sDATABASE", target))

	return ConnectionParameters{srcHost, srcPort, srcUser, srcPass, srcDb}
}

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
	flag.BoolVar(&tryOnly, "test", false, "Same as --try")

	flag.BoolVar(&purgeOnly, "purge", false, "Only purge destination tables and exit, no data will be inserted")

	flag.StringVar(&configFile, "config", "config.json", "Configuration file to use")

	// Local usage only
	debug := flag.Bool("debug", false, "run in loglevel DEBUG")
	flagVersion := flag.Bool("version", false, "display version")

	// Parse flags
	flag.Parse()

	if *flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}

// IsTerminal returns true if the file descriptor is connected to a terminal.
func IsTerminal(f *os.File) bool {
	// Get file descriptor information
	fileInfo, err := f.Stat()
	if err != nil {
		return false
	}

	// Check if the file descriptor mode indicates it's a terminal
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
