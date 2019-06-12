package helper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// loggers
var _stdout = log.New(os.Stdout, "", 0)
var _stderr = log.New(os.Stderr, "", 0)

// Constants
const (
	DefaultSessionKey      = "tmux"
	DefaultWindowName      = "new-window"
	DefaultSplitPercentage = 50
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
func IsSessionCreated(sessionName string, isVerbose bool) bool {
	args := []string{
		"has-session",
		"-t",
		sessionName,
	}

	if isVerbose {
		_stdout.Printf("[verbose] checking if session is created with command: tmux %s\n", strings.Join(args, " "))
	}

	if _, err := exec.Command("tmux", args...).CombinedOutput(); err == nil {
		return true
	}

	return false
}

// ListSessions list running sessions
func ListSessions(isVerbose bool) []string {
	args := []string{
		"ls",
	}

	if isVerbose {
		_stdout.Printf("[verbose] list running sessions with command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err == nil {
		return strings.Split(strings.TrimSpace(string(output)), "\n")
	}

	return []string{}
}

// StartSession starts a session
func (t *TmuxHelper) StartSession(sessionName string) bool {
	// check if 'tmux' is installed on the machine,
	if _, err := exec.LookPath("tmux"); err == nil {
		t.SessionName = sessionName

		return true
	}

	_stderr.Fatalf("* error finding 'tmux' on your system, please install it first.\n")

	return false
}

// CreateWindow creates a window
func (t *TmuxHelper) CreateWindow(windowName, directory, command string) bool {
	if t.SessionName != "" && IsSessionCreated(t.SessionName, t.Verbose) {
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
			_stdout.Printf("[verbose] creating a new window with command: tmux %s\n", strings.Join(args, " "))
		}

		output, err := exec.Command("tmux", args...).CombinedOutput()
		if err == nil {
			_stdout.Printf("> created a new window named: %s\n", windowName)

			if command != "" {
				t.Command(windowName, "", command)
			}

			return true
		}

		_stderr.Printf("* error creating a window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
	} else {
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
			_stdout.Printf("[verbose] creating a new session with command: tmux %s\n", strings.Join(args, " "))
		}

		output, err := exec.Command("tmux", args...).CombinedOutput()
		if err == nil {
			_stdout.Printf("> created a new session named: %s\n", t.SessionName)

			if command != "" {
				t.Command(windowName, "", command)
			}

			return true
		}

		_stderr.Printf("* error creating a new session: %s (%s)\n", t.SessionName, strings.TrimSpace(string(output)))
	}

	return false
}

// Command executes a command on a window/pane
func (t *TmuxHelper) Command(windowName string, paneName, command string) bool {
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

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		_stderr.Printf("* error executing command: %s (%s)\n", command, strings.TrimSpace(string(output)))
		return false
	}

	return true
}

// FocusWindow focuses on a window
func (t *TmuxHelper) FocusWindow(windowName string) bool {
	target := fmt.Sprintf("%s:%s", t.SessionName, windowName)
	args := []string{
		"select-window",
		"-t",
		target,
	}

	if t.Verbose {
		_stdout.Printf("[verbose] focusing window with command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		_stderr.Printf("* error focusing window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
		return false
	}

	return true
}

// FocusPane focuses on a pane
func (t *TmuxHelper) FocusPane(paneName string) bool {
	args := []string{
		"select-pane",
		"-t",
		paneName,
	}

	if t.Verbose {
		_stdout.Printf("[verbose] focusing pane with command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		_stderr.Printf("* error focusing pane: %s (%s)\n", paneName, strings.TrimSpace(string(output)))
		return false
	}

	return true
}

// SplitWindow splits a window
func (t *TmuxHelper) SplitWindow(windowName, directory string, options map[string]string) bool {
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
		_stdout.Printf("[verbose] splitting window with command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		_stderr.Printf("* error splitting window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
		return false
	}

	return true
}

// Attach attaches to a session
func (t *TmuxHelper) Attach() bool {
	command := []string{
		"tmux",
		"attach",
		"-t",
		t.SessionName,
	}

	path, err := exec.LookPath("tmux")
	if err == nil {
		if t.Verbose {
			_stdout.Printf("[verbose] attaching to a session with command: %s\n", strings.Join(command, " "))
		}

		err = syscall.Exec(path, command, syscall.Environ())
		if err == nil {
			return true
		}
	}

	_stderr.Printf("* error attaching to session: %s (%s)\n", t.SessionName, err)

	return false
}
