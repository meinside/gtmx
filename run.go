package main

import (
	"log"
	"os"
	"strings"

	"github.com/meinside/gtmx/config"
	"github.com/meinside/gtmx/tmux"

	"github.com/meinside/version-go"
)

// loggers
var (
	_stdout = log.New(os.Stdout, "", 0)
	_stderr = log.New(os.Stderr, "", 0)
)

// check for the existence of given param names
func paramExists(params []string, shortParam string, longParam string) bool {
	for _, param := range params {
		if param == shortParam || param == longParam {
			return true
		}
	}
	return false
}

// print usage and exits
func printUsageAndExit() {
	_stdout.Printf(`
> usage

# print this help message

$ gtmx -h
$ gtmx --help


# show sample config

$ gtmx -g
$ gtmx --gen-config


# show verbose messages (tmux commands)

$ gtmx -v
$ gtmx --verbose


# print version

$ gtmx -V
$ gtmx --version


# list predefined or running sessions

$ gtmx -l
$ gtmx --list


# kill this session

$ gtmx -q
$ gtmx --quit


# start, resume, or switch to a (predefined) session with/without the given key

$ gtmx [SESSION_KEY]
`)

	os.Exit(0)
}

// print version string and exits
func printVersionAndExit() {
	_stdout.Printf("%s", version.Minimum())

	os.Exit(0)
}

// print sample config file and exits
func printConfigAndExit() {
	sample := config.GetSampleConfigAsJSON()

	_stdout.Printf("/* sample config file (save it as $XDG_CONFIG_HOME/%s/%s) */\n\n", config.ApplicationName, config.ConfigFilename)

	_stdout.Println(sample)

	os.Exit(0)
}

// run with given parameters
func runWithParams(params []string, isVerbose bool) {
	var sessionKey string
	for _, param := range params {
		if !strings.HasPrefix(param, "-") {
			sessionKey = param
			continue
		}
	}
	if sessionKey == "" {
		var err error
		sessionKey, err = tmux.GetDefaultSessionKey()
		if err != nil {
			_stderr.Printf("* %s\n", err)
		}
	}

	errors := tmux.ConfigureAndAttachToSession(sessionKey, isVerbose)

	if len(errors) > 0 {
		for _, err := range errors {
			_stderr.Printf("%s\n", err)
		}
	}
}

// print sessions and exits
func printSessionsAndExit(isVerbose bool) {
	_stdout.Println()

	// list predefined sessions
	if confs := config.ReadAll(); len(confs) > 0 {
		_stdout.Printf("> all predefined sessions:\n")

		for name, conf := range confs {
			if conf.Description != nil {
				_stdout.Printf(" - %s: %s (%s)\n", name, conf.Name, *conf.Description)
			} else {
				_stdout.Printf(" - %s: %s\n", name, conf.Name)
			}
		}
	} else {
		_stdout.Printf("> no predefined sessions.\n")
	}

	_stdout.Println()

	// list running sessions
	sessions, err := tmux.ListSessions(isVerbose)
	if len(sessions) > 0 {
		_stdout.Printf("> all running sessions:\n")

		for _, session := range sessions {
			_stdout.Printf(" - %s\n", session)
		}
	} else {
		_stdout.Printf("> no running sessions.\n")
	}

	if isVerbose && err != nil {
		_stderr.Printf("* %s\n", err)
	}

	os.Exit(0)
}

// kill this session
func killCurrentSession() {
	if !tmux.IsInSession() {
		_stderr.Printf("* not in a tmux session\n")
		return
	}

	session, err := tmux.GetCurrentSessionName()
	if err == nil {
		err = tmux.KillSession(session)
		if err == nil {
			return
		}

		_stderr.Printf("* failed to kill session '%s': %s\n", session, err)
		return
	}

	_stderr.Printf("* %s\n", err)
}
