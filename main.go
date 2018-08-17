package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/meinside/gtmx/config"
	"github.com/meinside/gtmx/helper"
)

func paramExists(params []string, shortParam string, longParam string) bool {
	for _, param := range params {
		if param == shortParam || param == longParam {
			return true
		}
	}
	return false
}

func printUsageAndExit() {
	fmt.Printf(`
> Usage

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


# start or resume (predefined) session

$ gtmx [SESSION_NAME]
`)

	os.Exit(0)
}

func printConfigAndExit() {
	sample := config.GetSampleConfigAsJSON()

	fmt.Printf("/* Sample config file (save it as ~/%s) */\n\n", config.ConfigFilename)

	fmt.Println(sample)

	os.Exit(0)
}

func getDefaultSessionName() string {
	// new session
	output, err := exec.Command("hostname", "-s").CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	fmt.Printf("* Cannot get hostname, session name defaults to '%s'\n", helper.DefaultSessionName)

	return helper.DefaultSessionName
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
	run(isVerbose, params)
}

func run(isVerbose bool, params []string) {
	var sessionName string
	for _, param := range params {
		if !strings.HasPrefix(param, "-") {
			sessionName = param
			continue
		}
	}
	if sessionName == "" {
		sessionName = getDefaultSessionName()
	}

	tmux := helper.NewHelper()
	tmux.Verbose = isVerbose

	configs := config.ReadAll()

	if session, ok := configs[sessionName]; ok {
		fmt.Printf("> Using predefined session: %s\n", sessionName)

		session.Name = config.ReplaceString(session.Name)

		if session.RootDir != "" {
			fmt.Printf("> Session root directory: %s\n", session.RootDir)

			_, err := os.Stat(session.RootDir)

			if os.IsNotExist(err) {
				fmt.Printf("* Directory does not exist: %s\n", session.RootDir)
			} else {
				// change directory to it,
				if err := os.Chdir(session.RootDir); err != nil {
					fmt.Printf("* Failed to change directory: %s\n", session.RootDir)
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
			fmt.Printf("> Resuming session: %s\n", sessionName)

			tmux.StartSession(sessionName)
		}
	} else {
		tmux.StartSession(sessionName)

		if !helper.IsSessionCreated(sessionName, tmux.Verbose) {
			fmt.Printf("> No matching predefined session, creating a new session: %s\n", sessionName)

			tmux.CreateWindow(helper.DefaultWindowName, session.RootDir, "")
		} else {
			fmt.Printf("> No matching predefined session, resuming session: %s\n", sessionName)
		}
	}

	//attach
	tmux.Attach()
}

func printSessionsAndExit(isVerbose bool) {
	// predefined sessions
	fmt.Println()
	if confs := config.ReadAll(); len(confs) > 0 {
		fmt.Printf("> All predefined sessions:\n")

		for name, conf := range confs {
			if len(conf.Description) > 0 {
				fmt.Printf(" - %s (%s)\n", name, conf.Description)
			} else {
				fmt.Printf(" - %s\n", name)
			}
		}
	} else {
		fmt.Printf("> No predefined sessions.\n")
	}

	// running sessions
	fmt.Println()
	if sessions := helper.ListSessions(isVerbose); len(sessions) > 0 {
		fmt.Printf("> All running sessions:\n")

		for _, session := range sessions {
			fmt.Printf(" - %s\n", session)
		}
	} else {
		fmt.Printf("> No running sessions.\n")
	}

	os.Exit(0)
}
