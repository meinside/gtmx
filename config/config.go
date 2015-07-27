package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	CONFIG_FILENAME = ".gtmx.json" // config file's name
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

/*
	Read all predefined session configs from file
*/
func ReadAll() map[string]SessionConfig {
	all := make(map[string]SessionConfig)

	if user, err := user.Current(); err != nil {
		fmt.Printf("* Failed to get current user (%s)\n", err)

		os.Exit(1)
	} else {
		configFilepath := fmt.Sprintf("%s/%s", user.HomeDir, CONFIG_FILENAME)

		// config file exists,
		if _, err := os.Stat(configFilepath); err == nil {
			if file, err := ioutil.ReadFile(configFilepath); err != nil {
				fmt.Printf("* Failed to read config file (%s)\n", err)

				os.Exit(1)
			} else {
				if err := json.Unmarshal(file, &all); err != nil {
					fmt.Printf("* Failed to parse config file (%s)\n", err)

					os.Exit(1)
				}
			}
		}
	}

	return all
}

/*
	Get sample config (for generating sample config file)
*/
func GetSampleConfig() map[string]SessionConfig {
	sample := make(map[string]SessionConfig)

	sample["rails"] = SessionConfig{
		SessionName: "%d",
		Windows: []WindowConfig{
			{
				Name: "console",
			},
			{
				Name:    "models",
				Command: "cd ./app/models; clear",
			},
			{
				Name:    "views",
				Command: "cd ./app/views; clear",
			},
			{
				Name:    "controllers",
				Command: "cd ./app/controllers; clear",
			},
			{
				Name:    "configs",
				Command: "cd ./config; clear",
			},
			{
				Name:    "server",
				Command: "",
				Split: SplitConfig{
					Vertical:   true,
					Percentage: 50,
					Panes: []PaneConfig{
						{
							Pane:    "1",
							Command: "rails server",
						},
						{
							Pane:    "2",
							Command: "rails console",
						},
					},
				},
			},
		},
		Focus: FocusConfig{
			Name: "server",
			Pane: "2",
		},
	}

	return sample
}

/*
	Get sample config as JSON
*/
func GetSampleConfigAsJSON() string {
	sample := GetSampleConfig()
	if b, err := json.MarshalIndent(sample, "", "    "); err == nil {
		return string(b)
	}
	return "{}"
}

/*
	Replace string with place holders

	'%d' => current directory's name
	'%h' => hostname of this machine
*/
func ReplaceString(str string) string {
	replaced := str

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
