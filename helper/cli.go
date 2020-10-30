package helper

import (
	"os"
	"strings"

	"github.com/meinside/gtmx/config"
)

// ParamExists checks for the existence of given param names
func ParamExists(params []string, shortParam string, longParam string) bool {
	for _, param := range params {
		if param == shortParam || param == longParam {
			return true
		}
	}
	return false
}

// PrintUsageAndExit prints usage and exits
func PrintUsageAndExit() {
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


# list predefined or running sessions

$ gtmx -l
$ gtmx --list


# kill this session

$ gtmx -q
$ gtmx --quit


# start or resume a (predefined) session with/without given key

$ gtmx [SESSION_KEY]
`)

	os.Exit(0)
}

// PrintConfigAndExit prints sample config file and exits
func PrintConfigAndExit() {
	sample := config.GetSampleConfigAsJSON()

	_stdout.Printf("/* sample config file (save it as ~/%s) */\n\n", config.ConfigFilename)

	_stdout.Println(sample)

	os.Exit(0)
}

// RunWithParams runs with given parameters
func RunWithParams(params []string, isVerbose bool) {
	var sessionKey string
	for _, param := range params {
		if !strings.HasPrefix(param, "-") {
			sessionKey = param
			continue
		}
	}
	if sessionKey == "" {
		var err error
		sessionKey, err = GetDefaultSessionKey()

		if err != nil {
			_stderr.Printf("* %s\n", err)
		}
	}

	errors := ConfigureAndAttachToSession(sessionKey, isVerbose)

	if len(errors) > 0 {
		for _, err := range errors {
			_stderr.Printf("%s\n", err)
		}
	}
}

// PrintSessionsAndExit prints sessions and exits
func PrintSessionsAndExit(isVerbose bool) {
	_stdout.Println()

	// list predefined sessions
	if confs := config.ReadAll(); len(confs) > 0 {
		_stdout.Printf("> all predefined sessions:\n")

		for name, conf := range confs {
			if len(conf.Description) > 0 {
				_stdout.Printf(" - %s: %s (%s)\n", name, conf.Name, conf.Description)
			} else {
				_stdout.Printf(" - %s: %s\n", name, conf.Name)
			}
		}
	} else {
		_stdout.Printf("> no predefined sessions.\n")
	}

	_stdout.Println()

	// list running sessions
	sessions, err := ListSessions(isVerbose)
	if len(sessions) > 0 {
		_stdout.Printf("> all running sessions:\n")

		for _, session := range sessions {
			_stdout.Printf(" - %s\n", session)
		}
	} else {
		_stdout.Printf("> no running sessions.\n")
	}

	if err != nil {
		_stderr.Printf("* %s\n", err)
	}

	os.Exit(0)
}

// KillCurrentSession kills this session
func KillCurrentSession() {
	if !isInSession() {
		_stderr.Printf("* not in a tmux session\n")
		return
	}

	session, err := GetCurrentSessionName()
	if err == nil {
		err = KillSession(session)
		if err == nil {
			return
		}

		_stderr.Printf("* error killing session: %s (%s)\n", session, err)
		return
	}

	_stderr.Printf("* %s\n", err)
}
