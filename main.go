package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	config "github.com/meinside/gtmx/config"
	tmux "github.com/meinside/gtmx/helper"
)

func paramExists(params []string, shortParam string, longParam string) bool {
	for _, param := range params {
		if param == shortParam || param == longParam {
			return true
		}
	}
	return false
}

func printUsage() {
	fmt.Printf(`> Usage

# show this help message

$ gtmx -h
$ gtmx --help


# show sample config

$ gtmx -g
$ gtmx --gen-config


# show verbose messages (tmux commands)

$ gtmx -v
$ gtmx --verbose


# start or resume (predefined) session

$ gtmx [SESSION_NAME]
`)

	os.Exit(0)
}

func generateConfig() {
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

	fmt.Printf("* Cannot get hostname, session name defaults to '%s'\n", tmux.DefaultSessionName)

	return tmux.DefaultSessionName
}

func main() {
	params := os.Args[1:]

	// check params
	if paramExists(params, "-h", "--help") {
		printUsage()
	} else if paramExists(params, "-g", "--gen-config") {
		generateConfig()
	}

	var sessionName string
	if len(params) > 0 {
		sessionName = params[0]
	} else {
		sessionName = getDefaultSessionName()
	}

	helper := tmux.NewHelper()

	// check if verbose
	if paramExists(params, "-v", "--verbose") {
		helper.Verbose = true
	}

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

		if !tmux.IsSessionCreated(session.Name, helper.Verbose) {
			helper.StartSession(session.Name)

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
				helper.CreateWindow(windowName, dir, windowCommand)

				// split panes
				if window.Split.Percentage > 0 {
					helper.SplitWindow(windowName, dir, map[string]string{
						"vertical":   strconv.FormatBool(window.Split.Vertical),
						"percentage": strconv.Itoa(window.Split.Percentage),
					})
					for _, pane := range window.Split.Panes {
						helper.Command(windowName, pane.Pane, config.ReplaceString(pane.Command))
					}
				}
			}

			// focus window/pane
			if session.Focus.Name != "" {
				focusedWindow := session.Focus.Name
				focusedPane := session.Focus.Pane

				if focusedWindow != "" {
					helper.FocusWindow(focusedWindow)
					if focusedPane != "" {
						helper.FocusPane(focusedPane)
					}
				}
			}
		} else {
			fmt.Printf("> Resuming session: %s\n", sessionName)

			helper.StartSession(sessionName)
		}
	} else {
		helper.StartSession(sessionName)

		if !tmux.IsSessionCreated(sessionName, helper.Verbose) {
			fmt.Printf("> No matching predefined session, creating a new session: %s\n", sessionName)

			helper.CreateWindow(tmux.DefaultWindowName, session.RootDir, "")
		} else {
			fmt.Printf("> No matching predefined session, resuming session: %s\n", sessionName)
		}
	}

	//attach
	helper.Attach()
}
