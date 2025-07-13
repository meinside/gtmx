// config/config.go

// Package config for configuration of gtmx
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tailscale/hujson"
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
	Description *string        `json:"description,omitempty"`
	RootDir     *string        `json:"root_dir,omitempty"`
	Windows     []WindowConfig `json:"windows,omitempty"`
	Focus       *FocusConfig   `json:"focus,omitempty"`
}

// WindowConfig is a struct for window's configuration
type WindowConfig struct {
	Name        string       `json:"name"`
	Dir         *string      `json:"dir,omitempty"`
	Command     *string      `json:"cmd,omitempty"`
	Panes       []PaneConfig `json:"panes,omitempty"`
	Synchronize bool         `json:"synchronize,omitempty"`
}

// PaneConfig is a struct for pane's configuration
type PaneConfig struct {
	Name    string  `json:"name"`
	Command *string `json:"cmd,omitempty"`
}

// FocusConfig is a struct for focus' configuration
type FocusConfig struct {
	Name       string `json:"name"`
	PaneNumber *int   `json:"pane,omitempty"`
}

// ReadAll reads all predefined session configs from file.
func ReadAll() map[string]SessionConfig {
	all := make(map[string]SessionConfig)

	// https://xdgbasedirectoryspecification.com
	configDir := os.Getenv("XDG_CONFIG_HOME")

	// if the value of the environment variable is unset, empty, or not an absolute path, use the default one
	if configDir == "" || configDir[0:1] != "/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			_stderr.Fatalf(
				"* failed to get home directory: %s\n",
				err,
			)
		} else {
			configDir = filepath.Join(homeDir, ".config", ApplicationName)
		}
	} else {
		configDir = filepath.Join(configDir, ApplicationName)
	}

	configFilepath := fmt.Sprintf("%s/%s", configDir, ConfigFilename)

	// config file exists,
	if _, err := os.Stat(configFilepath); err == nil {
		if bytes, err := os.ReadFile(configFilepath); err != nil {
			_stderr.Fatalf(
				"* failed to read config file: %s\n",
				err,
			)
		} else {
			if bytes, err := standardizeJSON(bytes); err != nil {
				_stderr.Fatalf(
					"* failed to standardize config file to JWCC JSON: %s\n",
					err,
				)
			} else {
				if err := json.Unmarshal(bytes, &all); err != nil {
					_stderr.Fatalf(
						"* failed to parse config file: %s\n",
						err,
					)
				}
			}
		}
	}

	return all
}

// ToPtr returns the pointer of given value.
func ToPtr[T any](v T) *T {
	return &v
}

// ReplaceString replaces all placeholders in a given string.
//
// '%d' => current directory's name
// '%p' => current directory's path
// '%h' => hostname of this machine
func ReplaceString(str string) string {
	replaced := str

	// '%d' => current directory's name
	if strings.Contains(replaced, "%d") {
		if dir, err := os.Getwd(); err == nil {
			replaced = strings.ReplaceAll(replaced, "%d", filepath.Base(dir))
		}
	}

	// '%p' => current directory's path
	if strings.Contains(replaced, "%p") {
		if dir, err := os.Getwd(); err == nil {
			replaced = strings.ReplaceAll(replaced, "%p", dir)
		}
	}

	// '%h' => host name
	if strings.Contains(replaced, "%h") {
		if output, err := exec.Command("hostname", "-s").CombinedOutput(); err == nil {
			replaced = strings.ReplaceAll(replaced, "%h", strings.TrimSpace(string(output)))
		}
	}

	return replaced
}

// standardize given JSON (JWCC) bytes
func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()

	return ast.Pack(), nil
}
