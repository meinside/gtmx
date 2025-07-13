// tmux/tmux.go

// Package tmux for running tmux
package tmux

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/meinside/gtmx/config"
)

// loggers
var _stdout = log.New(os.Stdout, "", 0)

// Constants
const (
	DefaultSessionKey = "tmux"
	DefaultWindowName = "new-window"

	TmuxCommand = "tmux"
)

// TmuxHelper is a helper for tmux tasks
type TmuxHelper struct {
	SessionName string
	Verbose     bool
}

// NewHelper creates a new tmux helper.
func NewHelper() *TmuxHelper {
	return &TmuxHelper{}
}

// IsSessionCreated checks if a session is created or not.
func IsSessionCreated(sessionName string, isVerbose bool) (bool, error) {
	args := []string{
		"has-session",
		"-t",
		sessionName,
	}

	if isVerbose {
		_stdout.Printf(
			"[verbose] checking if session is created with command: `tmux %s`\n",
			strings.Join(args, " "),
		)
	}

	_, err := RunTmuxWithArgs(args)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ListSessions lists running sessions.
func ListSessions(isVerbose bool) (sessionLines []string, err error) {
	args := []string{
		"ls",
	}

	if isVerbose {
		_stdout.Printf(
			"[verbose] list running sessions with command: `tmux %s`\n",
			strings.Join(args, " "),
		)
	}

	output, err := RunTmuxWithArgs(args)
	if err == nil {
		return strings.Split(output, "\n"), nil
	}

	return []string{}, err
}

// GetDefaultSessionKey returns the default session key.
func GetDefaultSessionKey() (string, error) {
	// get hostname
	output, err := RunCommandWithArgs("hostname", []string{"-s"})
	if err == nil {
		return output, nil
	}

	return DefaultSessionKey, fmt.Errorf(
		"cannot get hostname, session key defaults to `%s` (%w)",
		DefaultSessionKey,
		err,
	)
}

// SetSessionName starts a session by naming it.
func (t *TmuxHelper) SetSessionName(sessionName string) error {
	// check if 'tmux' is installed on the machine,
	_, err := exec.LookPath(TmuxCommand)
	if err == nil {
		t.SessionName = sessionName
		return nil
	}

	return fmt.Errorf("`tmux` not found")
}

// CreateWindow creates a window.
func (t *TmuxHelper) CreateWindow(windowName string, directory, command *string) error {
	if t.SessionName != "" {
		if created, _ := IsSessionCreated(t.SessionName, t.Verbose); created {
			args := []string{
				"new-window",
				"-t",
				t.SessionName,
			}

			if windowName != "" {
				args = append(args, "-n", windowName)
			}
			if directory != nil {
				args = append(args, "-c", expandDir(*directory))
			}

			if t.Verbose {
				_stdout.Printf(
					"[verbose] creating a new window with command: `tmux %s`\n",
					strings.Join(args, " "),
				)
			}

			output, err := RunTmuxWithArgs(args)
			if err == nil {
				if t.Verbose {
					_stdout.Printf(
						"[verbose] created a new window named: %s\n",
						windowName,
					)
				}

				if command != nil {
					_ = t.Command(windowName, nil, *command)
				}

				return nil
			}

			return fmt.Errorf(
				"error creating a new window: %s (%s)",
				windowName,
				output,
			)
		}
	}

	// create a new session
	args := []string{
		"new-session",
		"-s",
		t.SessionName,
		"-n",
		windowName,
		"-d", // NOTE: detached
	}

	if directory != nil {
		args = append(args, "-c", expandDir(*directory))
	}

	if t.Verbose {
		_stdout.Printf(
			"[verbose] creating a new session with command: `tmux %s`\n",
			strings.Join(args, " "),
		)
	}

	output, err := RunTmuxWithArgs(args)
	if err == nil {
		if t.Verbose {
			_stdout.Printf(
				"[verbose] created a new session named: %s\n",
				t.SessionName,
			)
		}

		if command != nil {
			err = t.Command(windowName, nil, *command)
		}

		return err
	}

	return fmt.Errorf("error creating a new session: %s (%s)",
		t.SessionName,
		output,
	)
}

// RunCommandWithArgs runs a command with given arguments and returns the output.
func RunCommandWithArgs(cmd string, args []string) (output string, err error) {
	var bytes []byte
	bytes, err = exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf(
			"error executing `%s` with args: %s (%s)",
			cmd,
			strings.Join(args, " "),
			strings.TrimSpace(string(output)),
		)
	}
	return strings.TrimSpace(string(bytes)), err
}

// RunTmuxWithArgs runs `tmux` command with given arguments and returns the output.
func RunTmuxWithArgs(args []string) (output string, err error) {
	return RunCommandWithArgs(TmuxCommand, args)
}

// Command executes a command on a given window/pane.
func (t *TmuxHelper) Command(windowName string, paneName *string, command string) error {
	target := fmt.Sprintf("%s:%s", t.SessionName, windowName)
	if paneName != nil {
		target = fmt.Sprintf("%s.%s", target, *paneName)
	}

	args := []string{
		"send-keys",
		"-t",
		target,
		command,
		"C-m",
	}

	if t.Verbose {
		_stdout.Printf("[verbose] executing command: tmux %s\n", strings.Join(args, " "))
	}

	output, err := RunTmuxWithArgs(args)
	if err != nil {
		err = fmt.Errorf(
			"error executing command `%s` for target: %s (%s)",
			command,
			target,
			output,
		)
	}

	return err
}

// FocusWindow focuses on a window
func (t *TmuxHelper) FocusWindow(windowName string) error {
	target := fmt.Sprintf("%s:%s", t.SessionName, windowName)
	args := []string{
		"select-window",
		"-t",
		target,
	}

	if t.Verbose {
		_stdout.Printf(
			"[verbose] focusing window with command: `tmux %s`\n",
			strings.Join(args, " "),
		)
	}

	output, err := RunTmuxWithArgs(args)
	if err != nil {
		err = fmt.Errorf(
			"error focusing window: %s (%s)",
			target,
			output,
		)
	}

	return err
}

// FocusPane focuses on a pane.
func (t *TmuxHelper) FocusPane(paneNumber int) error {
	args := []string{
		"select-pane",
		"-t",
		strconv.Itoa(paneNumber),
	}

	if t.Verbose {
		_stdout.Printf(
			"[verbose] focusing pane with command: `tmux %s`\n",
			strings.Join(args, " "),
		)
	}

	output, err := RunTmuxWithArgs(args)
	if err != nil {
		err = fmt.Errorf(
			"error focusing pane: %d (%s)",
			paneNumber,
			output,
		)
	}

	return err
}

// SplitWindowTiled splits a window with tiled layout.
func (t *TmuxHelper) SplitWindowTiled(windowName string, directory *string, paneName string, cmd *string) error {
	target := fmt.Sprintf("%s:%s", t.SessionName, windowName)
	args := []string{
		"split-window",
		"-v",
		"-t",
		target,
		"-F",
		paneName,
	}
	if directory != nil {
		args = append(args, "-c", expandDir(*directory))
	}

	if t.Verbose {
		_stdout.Printf(
			"[verbose] splitting window tiled with command: `tmux %s`\n",
			strings.Join(args, " "),
		)
	}

	// split window,
	output, err := RunTmuxWithArgs(args)
	if err != nil {
		return fmt.Errorf(
			"error splitting window: %s (%s)",
			target,
			output,
		)
	}

	args = []string{
		"select-layout",
		"-t",
		target,
		"tiled",
	}

	if t.Verbose {
		_stdout.Printf(
			"[verbose] setting tiled layout with command: `tmux %s`\n",
			strings.Join(args, " "),
		)
	}

	// set tiled layout,
	output, err = RunTmuxWithArgs(args)
	if err != nil {
		return fmt.Errorf(
			"error setting tiled layout for target: %s (%s)",
			target,
			output,
		)
	}

	// and run command
	if cmd != nil {
		if err = t.Command(windowName, nil, *cmd); err != nil {
			return fmt.Errorf(
				"error running command for target: %s (%w)",
				target,
				err,
			)
		}
	}

	return nil
}

// Attach attaches to a session.
func (t *TmuxHelper) Attach() error {
	command := []string{
		TmuxCommand,
		"attach",
		"-t",
		t.SessionName,
	}

	path, err := exec.LookPath(TmuxCommand)
	if err == nil {
		if t.Verbose {
			_stdout.Printf(
				"[verbose] attaching to a session with command: `%s`\n",
				strings.Join(command, " "),
			)
		}

		err = syscall.Exec(path, command, syscall.Environ())
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf(
		"error attaching to session: %s (%w)",
		t.SessionName,
		err,
	)
}

// ConfigureAndAttachToSession configures up a session (if needed) and attaches to it.
func ConfigureAndAttachToSession(sessionKey string, isVerbose bool) (errors []error) {
	tmux := NewHelper()
	tmux.Verbose = isVerbose

	configs := config.ReadAll()
	errors = []error{}

	if session, ok := configs[sessionKey]; ok {
		if tmux.Verbose {
			_stdout.Printf(
				"[verbose] using predefined session with key: %s\n",
				sessionKey,
			)
		}

		session.Name = config.ReplaceString(session.Name)

		if tmux.Verbose {
			_stdout.Printf(
				"[verbose] using session name: %s\n",
				session.Name,
			)
		}

		if session.RootDir != nil {
			rootDir := expandDir(*session.RootDir)

			if tmux.Verbose {
				_stdout.Printf(
					"[verbose] session root directory: %s\n",
					rootDir,
				)
			}

			_, err := os.Stat(rootDir)

			if os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf(
					"directory does not exist: %s",
					rootDir,
				))
			} else {
				// change directory to it,
				if err := os.Chdir(rootDir); err != nil {
					errors = append(errors, fmt.Errorf(
						"failed to change directory: %s",
						rootDir,
					))
				}
			}
		}

		created, _ := IsSessionCreated(session.Name, tmux.Verbose)
		if !created {
			if err := tmux.SetSessionName(session.Name); err != nil {
				errors = append(errors, err)
			}

			for _, window := range session.Windows {
				// window name
				windowName := config.ReplaceString(window.Name)

				// window command
				windowCommand := window.Command
				if windowCommand != nil {
					windowCommand = config.ToPtr(config.ReplaceString(*windowCommand))
				}

				// create window with given name and command
				var dir *string
				if window.Dir != nil {
					dir = config.ToPtr(config.ReplaceString(*window.Dir))
				} else {
					dir = session.RootDir
				}
				if err := tmux.CreateWindow(windowName, dir, windowCommand); err != nil {
					errors = append(errors, err)
				}

				// split panes
				for _, pane := range window.Panes {
					var cmd *string
					if pane.Command != nil {
						cmd = config.ToPtr(config.ReplaceString(*pane.Command))
					}
					if err := tmux.SplitWindowTiled(windowName, dir, pane.Name, cmd); err != nil {
						errors = append(errors, err)
					}
				}

				// synchronize inputs
				target := fmt.Sprintf("%s:%s", session.Name, windowName)
				if window.Synchronize {
					args := []string{
						"set-window-option",
						"-t",
						target,
						"synchronize-panes",
						"on",
					}
					output, err := RunTmuxWithArgs(args)
					if err != nil {
						errors = append(errors, fmt.Errorf(
							"error synchronizing inputs for target: %s (%s)",
							target,
							output,
						))
					}
				}
			}

			// focus window/pane
			if session.Focus != nil && session.Focus.Name != "" {
				focusedWindow := session.Focus.Name

				if focusedWindow != "" {
					if err := tmux.FocusWindow(focusedWindow); err != nil {
						errors = append(errors, err)
					}
					if focusedPane := session.Focus.PaneNumber; focusedPane != nil {
						if err := tmux.FocusPane(*focusedPane); err != nil {
							errors = append(errors, err)
						}
					}
				}
			}
		} else {
			if tmux.Verbose {
				_stdout.Printf(
					"[verbose] resuming/switching to session: %s\n",
					session.Name,
				)
			}

			if err := tmux.SetSessionName(session.Name); err != nil {
				errors = append(errors, err)
			}

			// if already in another session, try switching to it instead of attaching
			if IsInSession() {
				if currentSessionName, err := GetCurrentSessionName(); err == nil {
					if currentSessionName != session.Name {
						if err := SwitchSession(session.Name); err == nil {
							return errors
						} else {
							errors = append(errors, err)
						}
					}
				} else {
					errors = append(errors, err)
				}
			}
		}
	} else {
		// use session key as a session name
		sessionName := sessionKey

		if err := tmux.SetSessionName(sessionName); err != nil {
			errors = append(errors, err)
		}

		created, _ := IsSessionCreated(sessionName, tmux.Verbose)
		if !created {
			if tmux.Verbose {
				_stdout.Printf(
					"[verbose] no matching predefined session, creating a new session: %s\n",
					sessionName,
				)
			}

			_ = tmux.CreateWindow(DefaultWindowName, session.RootDir, nil)
		} else {
			if tmux.Verbose {
				_stdout.Printf(
					"[verbose] no matching predefined session, resuming/switching to session: %s\n",
					sessionName,
				)
			}

			// if already in another session, try switching to it instead of attaching
			if IsInSession() {
				if currentSessionName, err := GetCurrentSessionName(); err == nil {
					if currentSessionName != sessionName {
						if err := SwitchSession(sessionName); err == nil {
							return errors
						} else {
							errors = append(errors, err)
						}
					}
				} else {
					errors = append(errors, err)
				}
			}
		}
	}

	// attach
	_ = tmux.Attach()

	return errors
}

// IsInSession checks if current session is in tmux.
func IsInSession() bool {
	env := os.Getenv("TMUX")

	return strings.TrimSpace(env) != ""
}

// GetCurrentSessionName returns current session's name.
func GetCurrentSessionName() (string, error) {
	args := []string{
		"display-message",
		"-p",
		"#S",
	}

	output, err := RunTmuxWithArgs(args)
	if err == nil {
		return output, nil
	}

	return "", fmt.Errorf(
		"failed to get current session name: %w",
		err,
	)
}

// SwitchSession switches to an existing session.
func SwitchSession(name string) error {
	command := []string{
		TmuxCommand,
		"switch",
		"-t",
		name,
	}

	path, err := exec.LookPath(TmuxCommand)
	if err == nil {
		err = syscall.Exec(path, command, syscall.Environ())
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf(
		"error switching to session: %s (%w)",
		name,
		err,
	)
}

// KillSession kills a session with given name.
func KillSession(name string) error {
	command := []string{
		TmuxCommand,
		"kill-session",
		"-t",
		name,
	}

	path, err := exec.LookPath(TmuxCommand)
	if err == nil {
		err = syscall.Exec(path, command, syscall.Environ())
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf(
		"error killing session: %s (%w)",
		name,
		err,
	)
}

// expand given directory's path (`~` and environment variables)
func expandDir(dir string) (expanded string) {
	expanded = dir

	// FIXME: expand dir paths with prefix: `~`
	if strings.HasPrefix(dir, "~/") { // ~/some/path
		home, _ := os.UserHomeDir()
		expanded = filepath.Join(home, dir[2:])
	} else if strings.HasPrefix(dir, "~") { // ~someuser/some/path
		splitted := strings.Split(dir, "/")
		username := splitted[0][1:] // drop `~`
		dirs := splitted[1:]
		home, _ := os.UserHomeDir()
		splitted = strings.Split(home, "/")
		splitted = append(splitted[:len(splitted)-1], username) // build up home path
		expanded = strings.Join(append(splitted, dirs...), "/") // append dirs to home path
	}

	// expand environment variables
	expanded = os.ExpandEnv(expanded)

	return
}
