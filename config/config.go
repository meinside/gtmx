package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// logger
var _stderr = log.New(os.Stderr, "", 0)

// Constants
const (
	ConfigFilename = ".gtmx.json" // config file's name
)

// SessionConfig is a struct for session's configuration
type SessionConfig struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	RootDir     string         `json:"root_dir,omitempty"`
	Windows     []WindowConfig `json:"windows,omitempty"`
	Focus       FocusConfig    `json:"focus,omitempty"`
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
		_stderr.Fatalf("* failed to get current user (%s)\n", err)
	} else {
		configFilepath := fmt.Sprintf("%s/%s", user.HomeDir, ConfigFilename)

		// config file exists,
		if _, err := os.Stat(configFilepath); err == nil {
			if file, err := ioutil.ReadFile(configFilepath); err != nil {
				_stderr.Fatalf("* failed to read config file (%s)\n", err)
			} else {
				if err := json.Unmarshal(file, &all); err != nil {
					_stderr.Fatalf("* failed to parse config file (%s)\n", err)
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
	// NOTE: This session should be started in a rails project directory.
	sample["Rails Project Template"] = SessionConfig{
		Name:        "rails-%d", // name session with current directory name
		Description: "predefined session for rails applications",
		Windows: []WindowConfig{
			{
				Name: "console",
			},
			{
				Name: "models",
				Dir:  "%p/app/models/", // relative directory
			},
			{
				Name: "views",
				Dir:  "%p/app/views/", // relative directory
			},
			{
				Name: "controllers", // relative directory
				Dir:  "%p/app/controllers/",
			},
			{
				Name: "configs",
				Dir:  "%p/config/",
			},
			{
				Name: "server",
				Split: SplitConfig{ // split into two panes and run commands
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
			Name: "server", // focus on the 'server' window
			Pane: "2",      // and '2' pane
		},
	}

	// (example 2) for this project
	sample["gtmx development"] = SessionConfig{
		Name:        "gtmx",
		Description: "predefined session for gtmx development",
		RootDir:     "/home/pi/go/src/github.com/meinside/gtmx", // absolute root directory
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
				Dir:     "%p/config/", // relative directory
				Command: "ls",
			},
			{
				Name:    "helper",
				Dir:     "%p/helper/", // relative directory
				Command: "ls",
			},
		},
		Focus: FocusConfig{
			Name: "main", // focus on the 'main' window
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
