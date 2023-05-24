package main

import (
	"os"

	"github.com/meinside/gtmx/helper"
)

func main() {
	params := os.Args[1:]

	// check if verbose option is on
	var isVerbose = helper.ParamExists(params, "-v", "--verbose")

	// check params
	if helper.ParamExists(params, "-h", "--help") {
		helper.PrintUsageAndExit()
	} else if helper.ParamExists(params, "-g", "--gen-config") {
		helper.PrintConfigAndExit()
	} else if helper.ParamExists(params, "-V", "--version") {
		helper.PrintVersionAndExit()
	} else if helper.ParamExists(params, "-l", "--list") {
		helper.PrintSessionsAndExit(isVerbose)
	} else if helper.ParamExists(params, "-q", "--quit") {
		helper.KillCurrentSession()
	} else {
		helper.RunWithParams(params, isVerbose)
	}
}
