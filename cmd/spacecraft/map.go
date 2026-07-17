package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func mapCmd(args []string, spaceDir, _ string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map <new|add|rm|ls|show|next|archive> [...]")
		return 1
	}

	roadmapsDir := filepath.Join(spaceDir, "roadmaps")
	os.MkdirAll(roadmapsDir, 0755)

	switch args[0] {
	case "new":
		return mapNew(args[1:], roadmapsDir)
	case "add":
		return mapAdd(args[1:], roadmapsDir)
	case "rm":
		return mapRemove(args[1:], roadmapsDir)
	case "ls":
		return mapList(roadmapsDir)
	case "show":
		return mapShow(args[1:], roadmapsDir)
	case "next":
		return mapNext(args[1:], roadmapsDir)
	case "archive":
		return mapArchive(args[1:], roadmapsDir, spaceDir)
	default:
		fmt.Fprintf(os.Stderr, "spacecraft map: unknown subcommand %q\n", args[0])
		return 1
	}
}

func mapNew(args []string, dir string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map new <title> [--desc <text>]")
		return 1
	}
	title := args[0]
	desc := ""
	for i := 1; i < len(args)-1; i++ {
		if args[i] == "--desc" {
			desc = args[i+1]
		}
	}

	id := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	id = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, id)

	now := time.Now().UTC().Format(time.RFC3339)
	r := map[string]interface{}{
		"id":          id,
		"title":       title,
		"description": desc,
		"missions":    []interface{}{},
		"issues":      []interface{}{},
		"createdAt":   now,
		"updatedAt":   now,
	}

	return saveRoadmap(dir, id, r)
}

func mapAdd(args []string, dir string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map add <roadmap-id> <mission-id> [--desc <text>]")
		return 1
	}
	rid, mid := args[0], args[1]
	desc := mid
	for i := 2; i < len(args)-1; i++ {
		if args[i] == "--desc" {
			desc = args[i+1]
		}
	}

	r, err := loadRoadmap(dir, rid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}

	missions, _ := r["missions"].([]interface{})
	r["missions"] = append(missions, map[string]interface{}{
		"id":          mid,
		"description": desc,
	})
	r["updatedAt"] = time.Now().UTC().Format(time.RFC3339)

	return saveRoadmap(dir, rid, r)
}

func mapRemove(args []string, dir string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map rm <roadmap-id> <mission-id>")
		return 1
	}
	rid, mid := args[0], args[1]

	r, err := loadRoadmap(dir, rid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}

	missions, _ := r["missions"].([]interface{})
	var filtered []interface{}
	for _, m := range missions {
		entry, _ := m.(map[string]interface{})
		if entry["id"] != mid {
			filtered = append(filtered, m)
		}
	}
	r["missions"] = filtered
	r["updatedAt"] = time.Now().UTC().Format(time.RFC3339)

	return saveRoadmap(dir, rid, r)
}

func mapList(dir string) int {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		r, err := loadRoadmap(dir, id)
		if err != nil {
			continue
		}
		missions, _ := r["missions"].([]interface{})
		issues, _ := r["issues"].([]interface{})
		fmt.Printf("%-30s %s (%d missions, %d issues)\n",
			id, r["title"], len(missions), len(issues))
	}
	return 0
}

func mapShow(args []string, dir string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map show <roadmap-id>")
		return 1
	}
	r, err := loadRoadmap(dir, args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}
	out, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(out))
	return 0
}

func mapNext(args []string, dir string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map next <roadmap-id>")
		return 1
	}
	r, err := loadRoadmap(dir, args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}
	missions, _ := r["missions"].([]interface{})
	for _, m := range missions {
		entry, _ := m.(map[string]interface{})
		if _, err := os.Stat(filepath.Join(dir, "..", "archive", entry["id"].(string))); os.IsNotExist(err) {
			fmt.Printf("%s: %s\n", entry["id"], entry["description"])
			return 0
		}
	}
	fmt.Println("All missions shipped.")
	return 0
}

func mapArchive(args []string, dir string, spaceDir string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map archive <roadmap-id>")
		return 1
	}
	id := args[0]
	archiveDir := filepath.Join(spaceDir, "archive", "roadmaps")
	os.MkdirAll(archiveDir, 0755)

	src := filepath.Join(dir, id+".json")
	dst := filepath.Join(archiveDir, id+".json")
	if err := os.Rename(src, dst); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}
	fmt.Printf("Archived roadmap: %s\n", id)
	return 0
}

func loadRoadmap(dir, id string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filepath.Join(dir, id+".json"))
	if err != nil {
		return nil, fmt.Errorf("roadmap not found: %s", id)
	}
	var r map[string]interface{}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("invalid roadmap: %v", err)
	}
	return r, nil
}

func saveRoadmap(dir, id string, r map[string]interface{}) int {
	data, _ := json.MarshalIndent(r, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, id+".json"), append(data, '\n'), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}
	fmt.Printf("Roadmap saved: %s\n", id)
	return 0
}

// Ensure sort imports work
var _ = sort.Strings
