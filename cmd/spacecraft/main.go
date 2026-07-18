package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage())
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft:", err)
		os.Exit(1)
	}

	spaceDir := filepath.Join(cwd, ".space")
	mid := resolveMission(cwd)

	switch cmd {
	case "init":
		os.Exit(initCmd(spaceDir))
	case "new":
		os.Exit(newCmd(args, spaceDir))
	case "missions":
		os.Exit(missionsCmd(spaceDir))
	case "use":
		os.Exit(useCmd(args, spaceDir))
	case "current":
		os.Exit(currentCmd(spaceDir))
	case "resolve":
		os.Exit(resolveCmd(args, spaceDir, mid))
	case "status":
		os.Exit(statusCmd(spaceDir, mid))
	case "flow":
		os.Exit(flowCmd(spaceDir, mid))
	case "bind-branch":
		os.Exit(bindBranchCmd(args, spaceDir, cwd, mid))
	case "git-info":
		os.Exit(gitInfoCmd(cwd))
	case "git-suggest":
		os.Exit(gitSuggestCmd(args, mid))
	case "state", "set-state":
		os.Exit(stateCmd(args, spaceDir, mid))
	case "clarify-status":
		os.Exit(clarifyStatusCmd(args, spaceDir, mid))
	case "evi", "evidence":
		os.Exit(eviCmd(args, spaceDir, mid))
	case "val", "validate":
		os.Exit(valCmd(args, spaceDir))
	case "closeout-check", "ship-check":
		os.Exit(closeoutCmd(spaceDir, mid))
	case "archive":
		os.Exit(archiveCmd(args, spaceDir, mid))
	case "map", "roadmap":
		os.Exit(mapCmd(args, spaceDir, cwd))
	case "-h", "--help", "help":
		fmt.Print(usage())
	default:
		fmt.Fprintf(os.Stderr, "spacecraft: unknown command %q\nTry 'spacecraft help'\n", cmd)
		os.Exit(1)
	}
}

func usage() string {
	return `Spacecraft mission helper

Usage:
  spacecraft init
      Initialize .space/ mission state.
  spacecraft new <title>
      Create a new mission with generated ID.
  spacecraft missions
      List all missions.
  spacecraft use <number|id|title>
      Select the current mission.
  spacecraft current
      Print the current mission ID.
  spacecraft resolve [selector]
      Resolve a mission from a selector or the current branch.
  spacecraft status
      Show mission status.
  spacecraft flow
      Show workflow snapshot for the resolved mission.
  spacecraft bind-branch [selector]
      Bind the current git branch to a mission.
  spacecraft git-info
      Show git worktree status.
  spacecraft git-suggest [type] [slug]
      Suggest a branch name and commit conventions.
  spacecraft set-state [mission-id] <new-state>
      Set mission state (alias: state).
  spacecraft clarify-status <open|clear|deferred>
      Set clarification status for the resolved mission.
  spacecraft evidence [--mission <id>] <label> -- <command...>
      Capture evidence; propagates the command exit code (alias: evi).
	  spacecraft validate [--strict] [mission-id]
	      Validate mission artifacts and evidence (alias: val).
	      --strict also requires exitCode on every evidence entry and
	      matching exitCode 0 evidence for each done plan task.
	  spacecraft closeout-check
	      Check whether a mission is ready to close out (alias: ship-check).
	  spacecraft ship-check
	      Alias for closeout-check.
  spacecraft archive [selector]
      Archive a shipped mission.
  spacecraft roadmap <new|add|rm|ls|show|next|archive> [...]
      Manage roadmaps in .space/roadmaps/ (alias: map).
  spacecraft help
      Show this help.

States: active -> planned -> in_progress -> ready -> shipped (blocked from any active state).
`
}

func missionDir(spaceDir, id string) string { return filepath.Join(spaceDir, "missions", id) }

func resolveMission(cwd string) string {
	out, _ := runCmd(cwd, "git", "branch", "--show-current")
	branch := strings.TrimSpace(out)
	if strings.HasPrefix(branch, "feat/") {
		parts := strings.SplitN(branch, "/", 3)
		if len(parts) >= 2 && strings.HasPrefix(parts[1], "M") {
			return parts[1]
		}
	}
	return ""
}
