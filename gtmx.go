package main

import (
	"fmt"
	config "github.com/meinside/gtmx/config"
	tmux "github.com/meinside/gtmx/helper"
	"os"
	"os/exec"
	"strconv"
	"strings"
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


# start or resume (predefined) session

$ gtmx [SESSION_NAME]
`)

	os.Exit(0)
}

func generateConfig() {
	sample := config.GetSampleConfigAsJSON()

	fmt.Printf("> Sample config (save it as ~/%s)\n\n", config.ConfigFilename)

	fmt.Println(sample)

	os.Exit(0)
}

func getDefaultSessionName() string {
	// new session
	if output, err := exec.Command("hostname", "-s").CombinedOutput(); err == nil {
		return strings.TrimSpace(string(output))
	} else {
		fmt.Printf("* Cannot get hostname, session name defaults to '%s'\n", tmux.DefaultSessionName)
		return tmux.DefaultSessionName
	}
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

	configs := config.ReadAll()
	if selected, ok := configs[sessionName]; ok {
		fmt.Printf("> Using predefined session: %s\n", sessionName)

		selected.SessionName = config.ReplaceString(selected.SessionName)

		if !tmux.IsSessionCreated(selected.SessionName) {
			helper.StartSession(selected.SessionName)

			for _, window := range selected.Windows {
				// window name
				windowName := config.ReplaceString(window.Name)

				// window command
				windowCommand := config.ReplaceString(window.Command)

				// create window with given name and command
				helper.CreateWindow(windowName, windowCommand)

				// split panes
				if window.Split.Percentage > 0 {
					helper.SplitWindow(windowName, map[string]string{
						"vertical":   strconv.FormatBool(window.Split.Vertical),
						"percentage": strconv.Itoa(window.Split.Percentage),
					})
					for _, pane := range window.Split.Panes {
						helper.Command(windowName, pane.Pane, config.ReplaceString(pane.Command))
					}
				}
			}

			// focus window/pane
			if selected.Focus.Name != "" {
				focusedWindow := selected.Focus.Name
				focusedPane := selected.Focus.Pane

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

		if !tmux.IsSessionCreated(sessionName) {
			fmt.Printf("> No matching predefined session, creating a new session: %s\n", sessionName)

			helper.CreateWindow(tmux.DefaultWindowName, "")
		} else {
			fmt.Printf("> No matching predefined session, resuming session: %s\n", sessionName)
		}
	}

	//attach
	helper.Attach()
}
