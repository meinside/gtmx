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

// Constants
const (
	ConfigFilename = ".gtmx.json" // config file's name
)

// SessionConfig is a struct for session's configuration
type SessionConfig struct {
	Name    string         `json:"name"`
	RootDir string         `json:"root_dir,omitempty"`
	Windows []WindowConfig `json:"windows,omitempty"`
	Focus   FocusConfig    `json:"focus,omitempty"`
}

// WindowConfig is a struct for window's configuration
type WindowConfig struct {
	Name    string      `json:"name"`
	Dir     string      `json:"dir,omitempty"`
	Command string      `json:"cmd,omitempty"`
	Split   SplitConfig `json:"split,omitempty"`
}

// SplitConfig is a struct for split's configuration
type SplitConfig struct {
	Vertical   bool         `json:"vertical,omitempty"`
	Percentage int          `json:"percentage,omitempty"`
	Panes      []PaneConfig `json:"panes,omitempty"`
}

// PaneConfig is a struct for pane's configuration
type PaneConfig struct {
	Pane    string `json:"pane"`
	Command string `json:"cmd,omitempty"`
}

// FocusConfig is a struct for focus' configuration
type FocusConfig struct {
	Name string `json:"name"`
	Pane string `json:"pane,omiempty"`
}

// ReadAll reads all predefined session configs from file
func ReadAll() map[string]SessionConfig {
	all := make(map[string]SessionConfig)

	if user, err := user.Current(); err != nil {
		fmt.Printf("* Failed to get current user (%s)\n", err)

		os.Exit(1)
	} else {
		configFilepath := fmt.Sprintf("%s/%s", user.HomeDir, ConfigFilename)

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

// getSampleConfig gets sample config (for generating sample config file)
func getSampleConfig() map[string]SessionConfig {
	sample := make(map[string]SessionConfig)

	// (example 1) for rails application
	sample["rails"] = SessionConfig{
		Name: "rails-%d",
		Windows: []WindowConfig{
			{
				Name: "console",
			},
			{
				Name: "models",
				Dir:  "%p/app/models/",
			},
			{
				Name: "views",
				Dir:  "%p/app/views/",
			},
			{
				Name: "controllers",
				Dir:  "%p/app/controllers/",
			},
			{
				Name: "configs",
				Dir:  "%p/config/",
			},
			{
				Name: "server",
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

	// (example 2) for this project
	sample["gtmx"] = SessionConfig{
		Name:    "gtmx",
		RootDir: "/home/pi/go/src/github.com/meinside/gtmx",
		Windows: []WindowConfig{
			{
				Name:    "git",
				Command: "git status",
			},
			{
				Name: "main",
			},
			{
				Name:    "config",
				Dir:     "%p/config/",
				Command: "ls",
			},
			{
				Name:    "helper",
				Dir:     "%p/helper/",
				Command: "ls",
			},
		},
		Focus: FocusConfig{
			Name: "main",
		},
	}

	return sample
}

// GetSampleConfigAsJSON generates a sample config as JSON string
func GetSampleConfigAsJSON() string {
	sample := getSampleConfig()
	if b, err := json.MarshalIndent(sample, "", "  "); err == nil {
		return string(b)
	}
	return "{}"
}

// ReplaceString replaces a string with place holders
//
// '%d' => current directory's name
// '%p' => current directory's path
// '%h' => hostname of this machine
func ReplaceString(str string) string {
	replaced := str

	// '%d' => current directory's name
	if strings.Contains(replaced, "%d") {
		if dir, err := os.Getwd(); err == nil {
			replaced = strings.Replace(replaced, "%d", filepath.Base(dir), -1)
		}
	}

	// '%p' => current directory's path
	if strings.Contains(replaced, "%p") {
		if dir, err := os.Getwd(); err == nil {
			replaced = strings.Replace(replaced, "%p", dir, -1)
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
