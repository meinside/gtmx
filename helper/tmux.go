package helper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Constants
const (
	DefaultSessionName     = "tmux"
	DefaultWindowName      = "new-window"
	DefaultSplitPercentage = "50"
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
		fmt.Printf("[verbose] Checking is session created with command: tmux %s\n", strings.Join(args, " "))
	}

	if _, err := exec.Command("tmux", args...).CombinedOutput(); err == nil {
		return true
	}

	return false
}

// StartSession starts a session
func (t *TmuxHelper) StartSession(sessionName string) bool {
	// check if 'tmux' is installed on the machine,
	if _, err := exec.LookPath("tmux"); err == nil {
		t.SessionName = sessionName

		return true
	}

	fmt.Printf("* Error finding 'tmux' on your system, please install it first.\n")
	os.Exit(1)

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
			fmt.Printf("[verbose] Creating a new window with command: tmux %s\n", strings.Join(args, " "))
		}

		output, err := exec.Command("tmux", args...).CombinedOutput()
		if err == nil {
			fmt.Printf("> Created a new window named: %s\n", windowName)

			if command != "" {
				t.Command(windowName, "", command)
			}

			return true
		}

		fmt.Printf("* Error creating a window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
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
			fmt.Printf("[verbose] Creating a new session with command: tmux %s\n", strings.Join(args, " "))
		}

		output, err := exec.Command("tmux", args...).CombinedOutput()
		if err == nil {
			fmt.Printf("> Created a new session named: %s\n", t.SessionName)

			if command != "" {
				t.Command(windowName, "", command)
			}

			return true
		}

		fmt.Printf("* Error creating a new session: %s (%s)\n", t.SessionName, strings.TrimSpace(string(output)))
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
		fmt.Printf("[verbose] Executing command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error executing command: %s (%s)\n", command, strings.TrimSpace(string(output)))
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
		fmt.Printf("[verbose] Focusing window with command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error focusing window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
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
		fmt.Printf("[verbose] Focusing pane with command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error focusing pane: %s (%s)\n", paneName, strings.TrimSpace(string(output)))
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
		if _, ok := options["vertical"]; ok {
			args = append(args, "-h")
		}

		// split percentage
		if option, ok := options["percentage"]; ok {
			args = append(args, []string{"-p", option}...)
		}

		// target pane
		if option, ok := options["pane"]; ok {
			target = fmt.Sprintf("%s.%s", target, option)
		}
	} else {
		args = append(args, "-v")
		args = append(args, []string{"-p", DefaultSplitPercentage}...)
	}

	args = append(args, []string{
		"-t",
		target,
	}...)

	if directory != "" {
		args = append(args, "-c", directory)
	}

	if t.Verbose {
		fmt.Printf("[verbose] Splitting window with command: tmux %s\n", strings.Join(args, " "))
	}

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error splitting window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
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
			fmt.Printf("[verbose] Attaching to a session with command: %s\n", strings.Join(command, " "))
		}

		err = syscall.Exec(path, command, syscall.Environ())
		if err == nil {
			return true
		}
	}

	fmt.Printf("* Error attaching to session: %s (%s)\n", t.SessionName, err)

	return false
}
