// main.go

package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

const (
	defaultUsage = `[OPTIONS...] [SESSION_NAME]`
)

func main() {
	// parse params,
	var p params
	parser := flags.NewParser(
		&p,
		flags.HelpFlag|flags.PassDoubleDash,
	)

	// set custom usage string
	parser.Usage = defaultUsage

	if remaining, err := parser.Parse(); err == nil {
		// check if multiple tasks were requested at a time
		if p.multipleTaskRequested() {
			os.Exit(
				printErrorBeforeExit(
					1,
					"Input error: multiple tasks were requested at a time.",
				),
			)
		}

		// run with params
		exit, err := run(p, remaining)

		if err != nil {
			os.Exit(
				printErrorBeforeExit(
					exit,
					"Error: %s",
					err,
				),
			)
		} else {
			os.Exit(exit)
		}
	} else {
		if e, ok := err.(*flags.Error); ok {
			helpExitCode := 0
			if e.Type != flags.ErrHelp {
				helpExitCode = 1

				printToStderrColored(
					color.FgHiRed,
					"Input error: %s",
					e.Error(),
				)
			}

			os.Exit(
				printHelpBeforeExit(
					helpExitCode,
					parser,
				),
			)
		}

		os.Exit(
			printErrorBeforeExit(
				1,
				"Failed to parse flags: %s",
				err,
			),
		)
	}

	// should not reach here
	os.Exit(
		printErrorBeforeExit(
			1,
			"Unhandled error.",
		),
	)
}
