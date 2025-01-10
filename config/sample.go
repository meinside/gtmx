package config

import "encoding/json"

// GetSampleConfig generates a sample config (for generating sample config file)
func GetSampleConfig() map[string]SessionConfig {
	sample := make(map[string]SessionConfig)

	// (example 1) for rails projects
	// NOTE: This session should be started in a rails project directory.
	sample["rails"] = SessionConfig{
		Name:        "rails-%d", // name session with current directory name
		Description: ToPtr("predefined session for rails projects"),
		Windows: []WindowConfig{
			{
				Name: "console",
			},
			{
				Name: "models",
				Dir:  ToPtr("%p/app/models/"), // relative directory
			},
			{
				Name: "views",
				Dir:  ToPtr("%p/app/views/"), // relative directory
			},
			{
				Name: "controllers", // relative directory
				Dir:  ToPtr("%p/app/controllers/"),
			},
			{
				Name: "configs",
				Dir:  ToPtr("%p/config/"),
			},
			{
				Name:    "server",
				Command: ToPtr("rails server"),
				Panes: []PaneConfig{
					{
						Name:    "console",
						Command: ToPtr("rails console"),
					},
				},
			},
		},
		Focus: &FocusConfig{
			Name:       "server", // focus on the 'server' window
			PaneNumber: ToPtr(2), // and pane 2 (console)
		},
	}

	// (example 2) for rust projects (created with rustup)
	// NOTE: This session should be started in a rust project directory.
	sample["rust"] = SessionConfig{
		Name:        "rust-%d",
		Description: ToPtr("predefined session for rust projects"),
		Windows: []WindowConfig{
			{
				Name:    "root",
				Command: ToPtr("git status"),
			},
			{
				Name:    "src",
				Dir:     ToPtr("%p/src/"), // relative directory
				Command: ToPtr("ls"),
			},
		},
		Focus: &FocusConfig{
			Name: "root", // focus on the 'root' window
		},
	}

	// (example 3) for clojure projects (created with lein)
	// NOTE: This session should be started in a clojure project directory.
	sample["clojure"] = SessionConfig{
		Name:        "clj-%d",
		Description: ToPtr("predefined session for clojure projects"),
		Windows: []WindowConfig{
			{
				Name:    "root",
				Command: ToPtr("git status"),
			},
			{
				Name:    "src",
				Dir:     ToPtr("%p/src/"), // relative directory
				Command: ToPtr("ls"),
			},
			{
				Name:    "test",
				Dir:     ToPtr("%p/test/"), // relative directory
				Command: ToPtr("ls"),
			},
			{
				Name:    "doc",
				Dir:     ToPtr("%p/doc/"), // relative directory
				Command: ToPtr("ls"),
			},
			{
				Name:    "repl",
				Dir:     ToPtr("%p/"), // relative directory
				Command: ToPtr("lein repl"),
			},
		},
		Focus: &FocusConfig{
			Name: "root", // focus on the 'root' window
		},
	}

	// (example 4) for managing multiple servers synchronously
	sample["multiple-servers"] = SessionConfig{
		Name:        "multiple servers",
		Description: ToPtr("for connecting to multiple servers and sending same commands all at once"),
		Windows: []WindowConfig{
			{
				Name:    "all-servers",
				Command: ToPtr("ssh user1@my-server-1 && exit"),
				Panes: []PaneConfig{
					{
						Name:    "server 2",
						Command: ToPtr("ssh user2@my-server-2 && exit"),
					},
					{
						Name:    "server 3",
						Command: ToPtr("ssh user3@my-server-3 && exit"),
					},
					{
						Name:    "server 4",
						Command: ToPtr("ssh user4@my-server-4 && exit"),
					},
					{
						Name:    "server 5",
						Command: ToPtr("ssh user5@my-server-5 && exit"),
					},
				},
				Synchronize: true, // synchronize inputs on all panes
			},
		},
		Focus: &FocusConfig{
			Name:       "all-servers",
			PaneNumber: ToPtr(5), // focus on the pane 5
		},
	}

	// (example 5) for developing this project
	sample["gtmx"] = SessionConfig{
		Name:        "gtmx-dev",
		Description: ToPtr("predefined session for gtmx development"),
		RootDir:     ToPtr("/home/ubuntu/go/src/github.com/meinside/gtmx"), // absolute root directory
		Windows: []WindowConfig{
			{
				Name:    "git",
				Command: ToPtr("git status"),
			},
			{
				Name: "main",
			},
			{
				Name:    "config",
				Dir:     ToPtr("%p/config/"), // relative directory
				Command: ToPtr("ls"),
			},
			{
				Name:    "tmux",
				Dir:     ToPtr("%p/tmux/"), // relative directory
				Command: ToPtr("ls"),
			},
		},
		Focus: &FocusConfig{
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
