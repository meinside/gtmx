package main

import (
	"os"
)

func main() {
	params := os.Args[1:]

	// check if verbose option is on
	var isVerbose = paramExists(params, "-v", "--verbose")

	// check params
	if paramExists(params, "-h", "--help") {
		printUsageAndExit()
	} else if paramExists(params, "-g", "--gen-config") {
		printConfigAndExit()
	} else if paramExists(params, "-V", "--version") {
		printVersionAndExit()
	} else if paramExists(params, "-l", "--list") {
		printSessionsAndExit(isVerbose)
	} else if paramExists(params, "-q", "--quit") {
		killCurrentSession()
	} else {
		runWithParams(params, isVerbose)
	}
}
