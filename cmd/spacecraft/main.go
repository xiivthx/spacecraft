package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
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
	case "evi":
		os.Exit(eviCmd(args, spaceDir, mid))
	case "val":
		os.Exit(valCmd(args, spaceDir))
	case "state":
		os.Exit(stateCmd(args, spaceDir))
	case "map":
		os.Exit(mapCmd(args, spaceDir, cwd))
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "spacecraft: unknown command %q\nTry 'spacecraft help'\n", cmd)
		os.Exit(1)
	}
}

func usage() {
	fmt.Print(`spacecraft — mission helper

  spacecraft evi <label> -- <command>
      Capture evidence. Runs command, saves output as JSONL.

  spacecraft val [mission-id]
      Validate mission artifacts exist and are well-formed.

  spacecraft state <mission-id> <new-state>
      Set mission state (active → planned → in_progress → ready → shipped).

  spacecraft map new <title> [--desc <text>]
  spacecraft map add <roadmap-id> <mission-id> [--desc <text>]
  spacecraft map rm <roadmap-id> <mission-id>
  spacecraft map ls
  spacecraft map show <roadmap-id>
  spacecraft map next <roadmap-id>
  spacecraft map archive <roadmap-id>
      Manage roadmaps in .space/roadmaps/.

  spacecraft help
`)
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
