package helper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	DEFAULT_SPLIT_PERCENTAGE = "50"
)

type TmuxHelper struct {
	SessionName string
}

func NewHelper() *TmuxHelper {
	return new(TmuxHelper)
}

func IsSessionCreated(sessionName string) bool {
	args := []string{
		"has-session",
		"-t",
		sessionName,
	}
	if _, err := exec.Command("tmux", args...).CombinedOutput(); err == nil {
		return true
	}
	return false
}

func IsWindowCreated(sessionName string, windowName string) bool {
	args := []string{
		"list-windows",
		"-t",
		sessionName,
		"-F",
		windowName,
	}
	if _, err := exec.Command("tmux", args...).CombinedOutput(); err == nil {
		return true
	}
	return false
}

func (t *TmuxHelper) StartSession(sessionName string) bool {
	// check if 'tmux' is installed on the machine,
	if _, err := exec.LookPath("tmux"); err == nil {
		if sessionName == "" {
			// new session
			if output, err := exec.Command("hostname", "-s").CombinedOutput(); err == nil {
				t.SessionName = strings.TrimSpace(string(output))
			}
		} else {
			// resume session
			t.SessionName = sessionName
		}

		return true
	} else {
		fmt.Printf("* Error finding 'tmux' on your system, please install it first.")
		os.Exit(1)
	}
	return false
}

func (t *TmuxHelper) CreateWindow(windowName string, command string) bool {
	if t.SessionName != "" && IsSessionCreated(t.SessionName) {
		args := []string{
			"new-window",
			"-t",
			t.SessionName,
		}
		if windowName != "" {
			args = append(args, "-n", windowName)
		}
		if output, err := exec.Command("tmux", args...).CombinedOutput(); err == nil {
			fmt.Printf("> Created a new window with name: %s\n", windowName)

			if command != "" {
				t.Command(windowName, "", command)
			}

			return true
		} else {
			fmt.Printf("* Error creating a window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
		}
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
		if output, err := exec.Command("tmux", args...).CombinedOutput(); err == nil {
			fmt.Printf("> Created a new session with name: %s\n", t.SessionName)
		} else {
			fmt.Printf("* Error creating a new session: %s (%s)\n", t.SessionName, strings.TrimSpace(string(output)))
		}
	}
	return false
}

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
	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error executing command: %s (%s)\n", command, strings.TrimSpace(string(output)))
		return false
	}
	return true
}

func (t *TmuxHelper) FocusWindow(windowName string) bool {
	target := fmt.Sprintf("%s:%s", t.SessionName, windowName)
	args := []string{
		"select-window",
		"-t",
		target,
	}
	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error focusing window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
		return false
	}
	return true
}

func (t *TmuxHelper) FocusPane(paneName string) bool {
	args := []string{
		"select-pane",
		"-t",
		paneName,
	}
	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error focusing pane: %s (%s)\n", paneName, strings.TrimSpace(string(output)))
		return false
	}
	return true
}

func (t *TmuxHelper) SplitWindow(windowName string, options map[string]string) bool {
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
		args = append(args, []string{"-p", DEFAULT_SPLIT_PERCENTAGE}...)
	}
	args = append(args, []string{
		"-t",
		target,
	}...)

	if output, err := exec.Command("tmux", args...).CombinedOutput(); err != nil {
		fmt.Printf("* Error splitting window: %s (%s)\n", windowName, strings.TrimSpace(string(output)))
		return false
	}
	return true
}

func (t *TmuxHelper) Attach() bool {
	command := []string{
		"tmux",
		"attach",
		"-t",
		t.SessionName,
	}

	if path, err := exec.LookPath("tmux"); err == nil {
		if err := syscall.Exec(path, command, syscall.Environ()); err != nil {
			fmt.Printf("* Error attaching to session: %s (%s)\n", t.SessionName, err)
		}
	}
	return true
}
