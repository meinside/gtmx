package helper

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/meinside/gtmx/config"
)

// Constants
const (
	DefaultSessionKey      = "tmux"
	DefaultWindowName      = "new-window"
	DefaultSplitPercentage = 50

	TmuxCommand = "tmux"
)

// TmuxHelper is a helper for tmux tasks
type TmuxHelper struct {
	SessionName string
	Verbose     bool
}

// NewHelper creates a new tmux helper
func NewHelper() *TmuxHelper {
	return new(TmuxHelper)
}

// IsSessionCreated checks if a session is created or not
func IsSessionCreated(sessionName string, isVerbose bool) (bool, error) {
	args := []string{
		"has-session",
		"-t",
		sessionName,
	}

	if isVerbose {
		_stdout.Printf("[verbose] checking if session is created with command: `tmux %s`\n", strings.Join(args, " "))
	}

	_, err := exec.Command(TmuxCommand, args...).CombinedOutput()
	if err != nil {
		return false, err
	}

	return true, nil
}

// ListSessions list running sessions
func ListSessions(isVerbose bool) (sessionLines []string, err error) {
	args := []string{
		"ls",
	}

	if isVerbose {
		_stdout.Printf("[verbose] list running sessions with command: `tmux %s`\n", strings.Join(args, " "))
	}

	output, err := exec.Command(TmuxCommand, args...).CombinedOutput()
	if err == nil {
		return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
	}

	return []string{}, err
}

// GetDefaultSessionKey returns the default session key
func GetDefaultSessionKey() (string, error) {
	// get hostname
	output, err := exec.Command("hostname", "-s").CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	return DefaultSessionKey, fmt.Errorf("cannot get hostname, session key defaults to `%s` (%s)", DefaultSessionKey, err)
}

// SetSessionName starts a session
func (t *TmuxHelper) SetSessionName(sessionName string) error {
	// check if 'tmux' is installed on the machine,
	_, err := exec.LookPath(TmuxCommand)
	if err == nil {
		t.SessionName = sessionName
		return nil
	}

	return fmt.Errorf("`tmux` not found")
}

// CreateWindow creates a window
func (t *TmuxHelper) CreateWindow(windowName, directory, command string) error {
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
			if directory != "" {
				args = append(args, "-c", directory)
			}

			if t.Verbose {
				_stdout.Printf("[verbose] creating a new window with command: `tmux %s`\n", strings.Join(args, " "))
			}

			output, err := exec.Command(TmuxCommand, args...).CombinedOutput()
			if err == nil {
				if t.Verbose {
					_stdout.Printf("[verbose] created a new window named: %s\n", windowName)
				}

				if command != "" {
					_ = t.Command(windowName, "", command)
				}

				return nil
			}

			return fmt.Errorf("error creating a new window: %s (%s)", windowName, strings.TrimSpace(string(output)))
		}
	}

	// create a new session
	args := []string{
		"new-session",
		"-s",
		t.SessionName,
		"-n",
		windowName,
		"-d",
	}

	if directory != "" {
		args = append(args, "-c", directory)
	}

	if t.Verbose {
		_stdout.Printf("[verbose] creating a new session with command: `tmux %s`\n", strings.Join(args, " "))
	}

	output, err := exec.Command(TmuxCommand, args...).CombinedOutput()
	if err == nil {
		if t.Verbose {
			_stdout.Printf("[verbose] created a new session named: %s\n", t.SessionName)
		}

		if command != "" {
			_ = t.Command(windowName, "", command)
		}

		return nil
	}

	return fmt.Errorf("error creating a new session: %s (%s)", t.SessionName, strings.TrimSpace(string(output)))
}

// Command executes a command on a window/pane
func (t *TmuxHelper) Command(windowName string, paneName, command string) error {
	target := fmt.Sprintf("%s:%s", t.SessionName, windowName)
	if paneName != "" {
		target = fmt.Sprintf("%s.%s", target, paneName)
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

	if output, err := exec.Command(TmuxCommand, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("error executing command: %s (%s)", command, strings.TrimSpace(string(output)))
	}

	return nil
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
		_stdout.Printf("[verbose] focusing window with command: `tmux %s`\n", strings.Join(args, " "))
	}

	if output, err := exec.Command(TmuxCommand, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("error focusing window: %s (%s)", windowName, strings.TrimSpace(string(output)))
	}

	return nil
}

// FocusPane focuses on a pane
func (t *TmuxHelper) FocusPane(paneName string) error {
	args := []string{
		"select-pane",
		"-t",
		paneName,
	}

	if t.Verbose {
		_stdout.Printf("[verbose] focusing pane with command: `tmux %s`\n", strings.Join(args, " "))
	}

	if output, err := exec.Command(TmuxCommand, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("error focusing pane: %s (%s)", paneName, strings.TrimSpace(string(output)))
	}

	return nil
}

// SplitWindow splits a window
func (t *TmuxHelper) SplitWindow(windowName, directory string, options map[string]string) error {
	args := []string{
		"split-window",
	}
	target := fmt.Sprintf("%s:%s", t.SessionName, windowName)

	if options != nil {
		// vertical split
		if v, exists := options["vertical"]; exists && v == "true" {
			args = append(args, "-h")
		}

		// split percentage
		if option, ok := options["percentage"]; ok {
			percent, err := strconv.Atoi(option)

			if err == nil {
				if percent < 10 || percent > 90 {
					percent = DefaultSplitPercentage
				}
			} else {
				percent = DefaultSplitPercentage
			}

			args = append(args, []string{"-p", strconv.Itoa(percent)}...)
		}

		// target pane
		if option, ok := options["pane"]; ok {
			target = fmt.Sprintf("%s.%s", target, option)
		}
	} else {
		args = append(args, "-v")
		args = append(args, []string{"-p", strconv.Itoa(DefaultSplitPercentage)}...)
	}

	args = append(args, []string{
		"-t",
		target,
	}...)

	if directory != "" {
		args = append(args, "-c", directory)
	}

	if t.Verbose {
		_stdout.Printf("[verbose] splitting window with command: `tmux %s`\n", strings.Join(args, " "))
	}

	if output, err := exec.Command(TmuxCommand, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("error splitting window: %s (%s)", windowName, strings.TrimSpace(string(output)))
	}

	return nil
}

// Attach attaches to a session
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
			_stdout.Printf("[verbose] attaching to a session with command: `%s`\n", strings.Join(command, " "))
		}

		err = syscall.Exec(path, command, syscall.Environ())
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("error attaching to session: %s (%s)", t.SessionName, err)
}

// ConfigureAndAttachToSession configures up a session (if needed) and attaches to it
func ConfigureAndAttachToSession(sessionKey string, isVerbose bool) (errors []error) {
	tmux := NewHelper()
	tmux.Verbose = isVerbose

	configs := config.ReadAll()
	errors = []error{}

	if session, ok := configs[sessionKey]; ok {
		if tmux.Verbose {
			_stdout.Printf("[verbose] using predefined session with key: %s\n", sessionKey)
		}

		session.Name = config.ReplaceString(session.Name)

		if tmux.Verbose {
			_stdout.Printf("[verbose] using session name: %s\n", session.Name)
		}

		if session.RootDir != "" {
			if tmux.Verbose {
				_stdout.Printf("[verbose] session root directory: %s\n", session.RootDir)
			}

			_, err := os.Stat(session.RootDir)

			if os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("directory does not exist: %s", session.RootDir))
			} else {
				// change directory to it,
				if err := os.Chdir(session.RootDir); err != nil {
					errors = append(errors, fmt.Errorf("failed to change directory: %s", session.RootDir))
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
				windowCommand := config.ReplaceString(window.Command)

				// create window with given name and command
				var dir string = window.Dir
				if dir == "" {
					dir = session.RootDir
				} else {
					dir = config.ReplaceString(dir)
				}
				if err := tmux.CreateWindow(windowName, dir, windowCommand); err != nil {
					errors = append(errors, err)
				}

				// split panes
				if window.Split.Percentage > 0 {
					if err := tmux.SplitWindow(windowName, dir, map[string]string{
						"vertical":   strconv.FormatBool(window.Split.Vertical),
						"percentage": strconv.Itoa(window.Split.Percentage),
					}); err != nil {
						errors = append(errors, err)
					}
					for _, pane := range window.Split.Panes {
						if err := tmux.Command(windowName, pane.Pane, config.ReplaceString(pane.Command)); err != nil {
							errors = append(errors, err)
						}
					}
				}
			}

			// focus window/pane
			if session.Focus.Name != "" {
				focusedWindow := session.Focus.Name
				focusedPane := session.Focus.Pane

				if focusedWindow != "" {
					if err := tmux.FocusWindow(focusedWindow); err != nil {
						errors = append(errors, err)
					}
					if focusedPane != "" {
						if err := tmux.FocusPane(focusedPane); err != nil {
							errors = append(errors, err)
						}
					}
				}
			}
		} else {
			if tmux.Verbose {
				_stdout.Printf("[verbose] resuming/switching to session: %s\n", session.Name)
			}

			if err := tmux.SetSessionName(session.Name); err != nil {
				errors = append(errors, err)
			}

			// if already in another session, try switching to it instead of attaching
			if isInSession() {
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
				_stdout.Printf("[verbose] no matching predefined session, creating a new session: %s\n", sessionName)
			}

			_ = tmux.CreateWindow(DefaultWindowName, session.RootDir, "")
		} else {
			if tmux.Verbose {
				_stdout.Printf("[verbose] no matching predefined session, resuming/switching to session: %s\n", sessionName)
			}

			// if already in another session, try switching to it instead of attaching
			if isInSession() {
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

	//attach
	_ = tmux.Attach()

	return errors
}

func isInSession() bool {
	env := os.Getenv("TMUX")

	return strings.TrimSpace(env) != ""
}

// GetCurrentSessionName returns current session's name
func GetCurrentSessionName() (string, error) {
	args := []string{
		"display-message",
		"-p",
		"#S",
	}

	output, err := exec.Command(TmuxCommand, args...).CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	return "", fmt.Errorf("failed to get current session name: %s", err)
}

// SwitchSession switches to an existing session
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

	return fmt.Errorf("error switching to session: %s (%s)", name, err)
}

// KillSession kills a session with given name
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

	return fmt.Errorf("error killing session: %s (%s)", name, err)
}
