// logging.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/jwalton/go-supportscolor"
)

// loggers
var (
	_stdout = log.New(os.Stdout, "", 0)
	_stderr = log.New(os.Stderr, "", 0)
)

// print help message before os.Exit()
func printHelpBeforeExit(
	code int,
	parser *flags.Parser,
) (exit int) {
	parser.WriteHelp(os.Stdout)

	return code
}

// print error before os.Exit()
func printErrorBeforeExit(
	code int,
	format string,
	a ...any,
) (exit int) {
	if code > 0 {
		printToStderrColored(
			color.FgHiRed,
			format,
			a...,
		)
	}

	return code
}

// print given string to stdout with color (if possible)
func printToStdoutColored(
	c color.Attribute,
	format string,
	a ...any,
) {
	formatted := fmt.Sprintf(format, a...)

	if supportscolor.Stdout().SupportsColor { // if color is supported,
		c := color.New(c)
		_, _ = c.Fprint(_stdout.Writer(), formatted)
	} else {
		fmt.Fprint(_stdout.Writer(), formatted)
	}
}

// print given string to stderr with color (if possible)
func printToStderrColored(
	c color.Attribute,
	format string,
	a ...any,
) {
	formatted := fmt.Sprintf(format, a...)

	if supportscolor.Stdout().SupportsColor { // if color is supported,
		c := color.New(c)
		_, _ = c.Fprint(_stderr.Writer(), formatted)
	} else {
		fmt.Fprint(_stderr.Writer(), formatted)
	}
}
