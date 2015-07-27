package main

import (
	"encoding/json"
	"fmt"
	tmux "github.com/meinside/gtmx/helper"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CONFIG_FILENAME = ".gtmx.json"
)

type SessionConfig struct {
	SessionName string         `json:"session_name"`
	Windows     []WindowConfig `json:"windows"`
	Focus       FocusConfig    `json:"focus"`
}
type WindowConfig struct {
	Name    string      `json:"name"`
	Command string      `json:"cmd"`
	Split   SplitConfig `json:"split"`
}
type SplitConfig struct {
	Vertical   bool         `json:"vertical"`
	Percentage int          `json:"percentage"`
	Panes      []PaneConfig `json:"panes"`
}
type PaneConfig struct {
	Pane    string `json:"pane"`
	Command string `json:"cmd"`
}
type FocusConfig struct {
	Name string `json:"name"`
	Pane string `json:"pane"`
}

func replaceString(sessionName string) string {
	replaced := sessionName

	// '%d' => current directory name
	if strings.Contains(replaced, "%d") {
		if dir, err := os.Getwd(); err == nil {
			replaced = strings.Replace(replaced, "%d", filepath.Base(dir), -1)
		}
	}

	// '%h' => host name
	if strings.Contains(replaced, "%h") {
		if output, err := exec.Command("hostname", "-s").CombinedOutput(); err == nil {
			replaced = strings.Replace(replaced, "%h", strings.TrimSpace(string(output)), 0)
		}
	}

	return replaced
}

func main() {
	if user, err := user.Current(); err != nil {
		fmt.Printf("* Failed to get current user (%s)\n", err)
	} else {
		helper := tmux.NewHelper()

		var sessionName string
		params := os.Args[1:]
		if len(params) > 0 {
			sessionName = params[0]
		} else {
			sessionName = ""
		}

		configFilepath := fmt.Sprintf("%s/%s", user.HomeDir, CONFIG_FILENAME)
		if file, err := ioutil.ReadFile(configFilepath); err == nil {
			var config map[string]SessionConfig
			if err := json.Unmarshal(file, &config); err != nil {
				fmt.Printf("* Failed to read config file (%s)\n", err)
			}

			if selectedConf, ok := config[sessionName]; ok {
				fmt.Printf("> Using predefined session: %s\n", sessionName)

				selectedConf.SessionName = replaceString(selectedConf.SessionName)

				if !tmux.IsSessionCreated(selectedConf.SessionName) {
					helper.StartSession(selectedConf.SessionName)

					for _, window := range selectedConf.Windows {
						// window name
						windowName := replaceString(window.Name)

						// window command
						windowCommand := replaceString(window.Command)

						// create window with given name and command
						helper.CreateWindow(windowName, windowCommand)

						// split panes
						if window.Split.Percentage > 0 {
							helper.SplitWindow(windowName, map[string]string{
								"vertical":   strconv.FormatBool(window.Split.Vertical),
								"percentage": strconv.Itoa(window.Split.Percentage),
							})
							for _, pane := range window.Split.Panes {
								helper.Command(windowName, pane.Pane, replaceString(pane.Command))
							}
						}
					}

					// focus window/pane
					if selectedConf.Focus.Name != "" {
						focusedWindow := selectedConf.Focus.Name
						focusedPane := selectedConf.Focus.Pane

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

					helper.CreateWindow(tmux.DEFAULT_WINDOW_NAME, "")
				} else {
					fmt.Printf("> No matching predefined session, resuming session: %s\n", sessionName)
				}
			}
		} else {
			helper.StartSession(sessionName)

			if !tmux.IsSessionCreated(sessionName) {
				fmt.Printf("> Not using predefined session, creating a new session: %s\n", sessionName)

				helper.CreateWindow(tmux.DEFAULT_WINDOW_NAME, "")
			} else {
				fmt.Printf("> Not using predefined session, resuming session: %s\n", sessionName)
			}
		}

		//attach
		helper.Attach()
	}
}
