package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// standard error
var _stderr = log.New(os.Stderr, "", 0)

// Constants
const (
	ApplicationName = "gtmx"
	ConfigFilename  = "config.json" // config file's name
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
	Pane string `json:"pane,omitempty"`
}

// ReadAll reads all predefined session configs from file
func ReadAll() map[string]SessionConfig {
	all := make(map[string]SessionConfig)

	// https://xdgbasedirectoryspecification.com
	configDir := os.Getenv("XDG_CONFIG_HOME")

	// If the value of the environment variable is unset, empty, or not an absolute path, use the default
	if configDir == "" || configDir[0:1] != "/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			_stderr.Fatalf("* failed to get home directory (%s)\n", err)
		} else {
			configDir = filepath.Join(homeDir, ".config", ApplicationName)
		}
	} else {
		configDir = filepath.Join(configDir, ApplicationName)
	}

	configFilepath := fmt.Sprintf("%s/%s", configDir, ConfigFilename)

	// config file exists,
	if _, err := os.Stat(configFilepath); err == nil {
		if file, err := os.ReadFile(configFilepath); err != nil {
			_stderr.Fatalf("* failed to read config file (%s)\n", err)
		} else {
			if err := json.Unmarshal(file, &all); err != nil {
				_stderr.Fatalf("* failed to parse config file (%s)\n", err)
			}
		}
	}

	return all
}

// GetSampleConfig generates a sample config (for generating sample config file)
func GetSampleConfig() map[string]SessionConfig {
	sample := make(map[string]SessionConfig)

	// (example 1) for rails projects
	// NOTE: This session should be started in a rails project directory.
	sample["rails"] = SessionConfig{
		Name:        "rails-%d", // name session with current directory name
		Description: "predefined session for rails projects",
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

	// (example 2) for rust projects (created with rustup)
	// NOTE: This session should be started in a rust project directory.
	sample["rust"] = SessionConfig{
		Name:        "rust-%d",
		Description: "predefined session for rust projects",
		Windows: []WindowConfig{
			{
				Name:    "root",
				Command: "git status",
			},
			{
				Name:    "src",
				Dir:     "%p/src/", // relative directory
				Command: "ls",
			},
		},
		Focus: FocusConfig{
			Name: "root", // focus on the 'root' window
		},
	}

	// (example 3) for clojure projects (created with lein)
	// NOTE: This session should be started in a clojure project directory.
	sample["clojure"] = SessionConfig{
		Name:        "clj-%d",
		Description: "predefined session for clojure projects",
		Windows: []WindowConfig{
			{
				Name:    "root",
				Command: "git status",
			},
			{
				Name:    "src",
				Dir:     "%p/src/", // relative directory
				Command: "ls",
			},
			{
				Name:    "test",
				Dir:     "%p/test/", // relative directory
				Command: "ls",
			},
			{
				Name:    "doc",
				Dir:     "%p/doc/", // relative directory
				Command: "ls",
			},
			{
				Name:    "repl",
				Dir:     "%p/", // relative directory
				Command: "lein repl",
			},
		},
		Focus: FocusConfig{
			Name: "root", // focus on the 'root' window
		},
	}

	// (example 4) for babashka scripts and conjure
	sample["bb"] = SessionConfig{
		Name:        "bb-%d",
		Description: "predefined session for babashka scripts and conjure with nrepl connection",
		Windows: []WindowConfig{
			{
				Name:    "nrepl",
				Dir:     "%p/", // relative directory
				Command: "echo `shuf -i 10000-50000 -n 1` > .nrepl-port && bb --nrepl-server `cat .nrepl-port`",
			},
			{
				Name:    "scripts",
				Dir:     "%p/", // relative directory
				Command: "ls",
			},
		},
		Focus: FocusConfig{
			Name: "scripts", // focus on the 'scripts' window
		},
	}

	// (example 5) for this project
	sample["gtmx"] = SessionConfig{
		Name:        "gtmx-dev",
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
	sample := GetSampleConfig()
	if b, err := json.MarshalIndent(sample, "", "  "); err == nil {
		return string(b)
	}
	return "{}"
}

// ReplaceString replaces all place holders in a given string
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
			replaced = strings.Replace(replaced, "%h", strings.TrimSpace(string(output)), -1)
		}
	}

	return replaced
}
