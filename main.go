package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/meinside/gtmx/config"
	"github.com/meinside/gtmx/helper"
)

// loggers
var _stdout = log.New(os.Stdout, "", 0)
var _stderr = log.New(os.Stderr, "", 0)

func paramExists(params []string, shortParam string, longParam string) bool {
	for _, param := range params {
		if param == shortParam || param == longParam {
			return true
		}
	}
	return false
}

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


# list predefined or running sessions

$ gtmx -l
$ gtmx --list


# start or resume a (predefined) session with/without given key

$ gtmx [SESSION_KEY]
`)

	os.Exit(0)
}

func printConfigAndExit() {
	sample := config.GetSampleConfigAsJSON()

	_stdout.Printf("/* sample config file (save it as ~/%s) */\n\n", config.ConfigFilename)

	_stdout.Println(sample)

	os.Exit(0)
}

func getDefaultSessionKey() string {
	// get hostname
	output, err := exec.Command("hostname", "-s").CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	_stderr.Printf("* cannot get hostname, session key defaults to '%s'\n", helper.DefaultSessionKey)

	return helper.DefaultSessionKey
}

func main() {
	params := os.Args[1:]

	// check if verbose option is on
	var isVerbose = paramExists(params, "-v", "--verbose")

	// check params
	if paramExists(params, "-h", "--help") {
		printUsageAndExit()
	} else if paramExists(params, "-g", "--gen-config") {
		printConfigAndExit()
	} else if paramExists(params, "-l", "--list") {
		printSessionsAndExit(isVerbose)
	}

	// run with params
	run(params, isVerbose)
}

func run(params []string, isVerbose bool) {
	var sessionKey string
	for _, param := range params {
		if !strings.HasPrefix(param, "-") {
			sessionKey = param
			continue
		}
	}
	if sessionKey == "" {
		sessionKey = getDefaultSessionKey()
	}

	tmux := helper.NewHelper()
	tmux.Verbose = isVerbose

	configs := config.ReadAll()

	if session, ok := configs[sessionKey]; ok {
		_stdout.Printf("> using predefined session with key: %s\n", sessionKey)

		session.Name = config.ReplaceString(session.Name)

		_stdout.Printf("> using session name: %s\n", session.Name)

		if session.RootDir != "" {
			_stdout.Printf("> session root directory: %s\n", session.RootDir)

			_, err := os.Stat(session.RootDir)

			if os.IsNotExist(err) {
				_stderr.Printf("* directory does not exist: %s\n", session.RootDir)
			} else {
				// change directory to it,
				if err := os.Chdir(session.RootDir); err != nil {
					_stderr.Printf("* failed to change directory: %s\n", session.RootDir)
				}
			}
		}

		if !helper.IsSessionCreated(session.Name, tmux.Verbose) {
			tmux.StartSession(session.Name)

			for _, window := range session.Windows {
				// window name
				windowName := config.ReplaceString(window.Name)

				// window command
				windowCommand := config.ReplaceString(window.Command)

				// create window with given name and command
				var dir string = window.Dir
				if dir == "" {
					dir = session.RootDir
				} else {
					dir = config.ReplaceString(dir)
				}
				tmux.CreateWindow(windowName, dir, windowCommand)

				// split panes
				if window.Split.Percentage > 0 {
					tmux.SplitWindow(windowName, dir, map[string]string{
						"vertical":   strconv.FormatBool(window.Split.Vertical),
						"percentage": strconv.Itoa(window.Split.Percentage),
					})
					for _, pane := range window.Split.Panes {
						tmux.Command(windowName, pane.Pane, config.ReplaceString(pane.Command))
					}
				}
			}

			// focus window/pane
			if session.Focus.Name != "" {
				focusedWindow := session.Focus.Name
				focusedPane := session.Focus.Pane

				if focusedWindow != "" {
					tmux.FocusWindow(focusedWindow)
					if focusedPane != "" {
						tmux.FocusPane(focusedPane)
					}
				}
			}
		} else {
			_stdout.Printf("> resuming session: %s\n", session.Name)

			tmux.StartSession(session.Name)
		}
	} else {
		// use session key as a session name
		sessionName := sessionKey

		tmux.StartSession(sessionName)

		if !helper.IsSessionCreated(sessionName, tmux.Verbose) {
			_stdout.Printf("> no matching predefined session, creating a new session: %s\n", sessionName)

			tmux.CreateWindow(helper.DefaultWindowName, session.RootDir, "")
		} else {
			_stdout.Printf("> no matching predefined session, resuming session: %s\n", sessionName)
		}
	}

	//attach
	tmux.Attach()
}

func printSessionsAndExit(isVerbose bool) {
	// predefined sessions
	_stdout.Println()
	if confs := config.ReadAll(); len(confs) > 0 {
		_stdout.Printf("> all predefined sessions:\n")

		for name, conf := range confs {
			if len(conf.Description) > 0 {
				_stdout.Printf(" - %s (%s)\n", name, conf.Description)
			} else {
				_stdout.Printf(" - %s\n", name)
			}
		}
	} else {
		_stdout.Printf("> no predefined sessions.\n")
	}

	// running sessions
	_stdout.Println()
	if sessions := helper.ListSessions(isVerbose); len(sessions) > 0 {
		_stdout.Printf("> all running sessions:\n")

		for _, session := range sessions {
			_stdout.Printf(" - %s\n", session)
		}
	} else {
		_stdout.Printf("> no running sessions.\n")
	}

	os.Exit(0)
}
