package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/meinside/gtmx/config"
	"github.com/meinside/gtmx/tmux"

	"github.com/meinside/version-go"
)

// run with params and arguments
func run(
	p params,
	remainingArgs []string,
) (exit int, err error) {
	isVerbose := p.Verbose

	// handle predefined tasks
	if p.PrintVersion {
		return printVersionAndExit()
	} else if p.GenerateConfig {
		return printConfigAndExit()
	} else if p.ListSessions {
		return printSessionsAndExit(isVerbose)
	} else if p.QuitCurrentSession {
		return killCurrentSession()
	}

	// fallback with remaining arguments
	return runWithArgs(remainingArgs, isVerbose)
}

// print version string and exit
func printVersionAndExit() (exit int, err error) {
	printToStdoutColored(
		color.FgHiGreen,
		"%s\n",
		version.Minimum(),
	)

	return 0, nil
}

// print sample config file and exit
func printConfigAndExit() (exit int, err error) {
	sample := config.GetSampleConfigAsJSON()

	printToStdoutColored(
		color.FgCyan,
		"/* sample config file (save it as $XDG_CONFIG_HOME/%s/%s) */\n",
		config.ApplicationName,
		config.ConfigFilename,
	)

	printToStdoutColored(
		color.FgHiCyan,
		"%s\n",
		sample,
	)

	return 0, nil
}

// print sessions and exit
func printSessionsAndExit(isVerbose bool) (code int, err error) {
	_stdout.Println()

	// list predefined sessions
	if confs := config.ReadAll(); len(confs) > 0 {
		printToStdoutColored(
			color.FgWhite,
			"> all predefined sessions:\n",
		)

		for name, conf := range confs {
			var line string

			if conf.Description != nil {
				line = fmt.Sprintf(
					" - %s: %s (%s)\n",
					escape(name),
					escape(conf.Name),
					escape(*conf.Description),
				)
			} else {
				line = fmt.Sprintf(
					" - %s: %s\n",
					escape(name),
					escape(conf.Name),
				)
			}

			printToStdoutColored(
				color.FgHiWhite,
				line,
			)
		}
	} else {
		printToStdoutColored(
			color.FgWhite,
			"> no predefined sessions.\n",
		)
	}

	_stdout.Println()

	// list running sessions
	var sessions []string
	sessions, err = tmux.ListSessions(isVerbose)
	if len(sessions) > 0 {
		printToStdoutColored(
			color.FgWhite,
			"> all running sessions:\n",
		)

		for _, session := range sessions {
			printToStdoutColored(
				color.FgHiWhite,
				" - %s\n",
				session,
			)
		}
	} else {
		printToStdoutColored(
			color.FgWhite,
			"> no running sessions.\n",
		)
	}

	if err != nil {
		if isVerbose {
			printToStderrColored(
				color.FgHiRed,
				"* %s\n",
				err,
			)
		}
	}

	return 0, nil
}

// kill this session
func killCurrentSession() (code int, err error) {
	if !tmux.IsInSession() {
		return 1, fmt.Errorf("not in a tmux session")
	}

	session, err := tmux.GetCurrentSessionName()
	if err == nil {
		err = tmux.KillSession(session)
		if err == nil {
			return
		}

		return 1, fmt.Errorf(
			"failed to kill session '%s': %s",
			session,
			err,
		)
	} else {
		err = fmt.Errorf(
			"failed to get current session for killing: %s",
			err,
		)
	}

	return 1, err
}

// run with given arguments
func runWithArgs(args []string, isVerbose bool) (exit int, err error) {
	// take the first session name
	var sessionKey string
	if len(args) > 0 {
		sessionKey = args[0]
	}

	// if there is no session name given, take the default one
	if sessionKey == "" {
		if sessionKey, err = tmux.GetDefaultSessionKey(); err != nil {
			return 1, fmt.Errorf(
				"failed to get the default session key: %s",
				err,
			)
		}
	}

	// configure and attach to given session name
	if errs := tmux.ConfigureAndAttachToSession(sessionKey, isVerbose); len(errs) > 0 {
		return 1, errors.Join(errs...)
	}

	return 0, nil
}
