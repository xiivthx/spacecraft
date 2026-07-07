package main

import (
	"fmt"
	"os"
)

func usage() string {
	return `Spacecraft local mission helper

Usage:
  spacecraft init
  spacecraft new <title>
  spacecraft current
  spacecraft resolve [selector] [--json]
  spacecraft missions
  spacecraft use <number|id|title>
  spacecraft bind-branch [selector]
  spacecraft status
  spacecraft flow [--json]
  spacecraft git-info
  spacecraft git-suggest [type] [slug]
  spacecraft set-state <state>
  spacecraft clarify-status <open|clear|deferred>
  spacecraft evidence <label> -- <command...>
  spacecraft validate
  spacecraft closeout-check
  spacecraft archive [selector]
`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage())
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "init":
		initSpacecraft(false)
	case "new":
		createMission(stringsJoin(args, " "))
	case "current":
		printCurrent()
	case "resolve":
		printResolvedMission(args)
	case "missions":
		printMissions()
	case "use":
		useMission(args)
	case "bind-branch":
		bindBranch(args)
	case "status":
		printStatus()
	case "flow":
		printWorkflow(args)
	case "git-info":
		printGitInfo()
	case "git-suggest":
		printGitSuggestion(args)
	case "set-state":
		if len(args) == 0 {
			fail("Missing state.\nAllowed states: draft, specified, planned, implementing, verifying, reviewing, ready, shipped, blocked")
		}
		setState(args[0])
	case "clarify-status":
		if len(args) == 0 {
			fail("Missing clarification status.\nAllowed statuses: open, clear, deferred")
		}
		setClarificationStatus(args[0])
	case "evidence":
		recordEvidence(args)
	case "validate":
		validateMission()
	case "closeout-check":
		releaseCloseoutCheck()
	case "archive":
		archiveMission(args)
	case "-h", "--help", "help":
		fmt.Print(usage())
	default:
		fail(fmt.Sprintf("Unknown command %q.\n\n%s", command, usage()))
	}
}

func stringsJoin(args []string, sep string) string {
	res := ""
	for i, a := range args {
		if i > 0 {
			res += sep
		}
		res += a
	}
	return res
}

// Stubs to allow compilation before all files are written
func releaseCloseoutCheck()                {}
func archiveMission(args []string)         {}
