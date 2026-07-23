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
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map <new|add|rm|ls|show|next|use|current|archive> [...]")
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
	case "use":
		return mapUse(args[1:], spaceDir, roadmapsDir)
	case "current":
		return mapCurrent(spaceDir)
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
	spaceDir := filepath.Join(dir, "..")
	if id, desc, state, ok := nextIncompleteOnRoadmap(spaceDir, r); ok {
		fmt.Printf("%s: %s (state=%s)\n", id, desc, state)
		return 0
	}
	fmt.Println("All missions complete.")
	return 0
}

// nextIncompleteOnRoadmap returns the first incomplete (non-ready, non-shipped,
// non-archived) mission on a roadmap.
func nextIncompleteOnRoadmap(spaceDir string, r map[string]interface{}) (id, desc, state string, ok bool) {
	missions, _ := r["missions"].([]interface{})
	for _, m := range missions {
		entry, _ := m.(map[string]interface{})
		mid, _ := entry["id"].(string)
		d, _ := entry["description"].(string)
		if mid == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(spaceDir, "archive", mid)); err == nil {
			continue
		}
		st := missionIncompleteState(filepath.Join(spaceDir, "missions", mid))
		if st == "" {
			continue
		}
		return mid, d, st, true
	}
	return "", "", "", false
}

func roadmapContainsMission(r map[string]interface{}, missionID string) bool {
	missions, _ := r["missions"].([]interface{})
	for _, m := range missions {
		entry, _ := m.(map[string]interface{})
		if id, _ := entry["id"].(string); id == missionID {
			return true
		}
	}
	return false
}

// findRoadmapForArchivedMission prefers current-roadmap when it lists the
// archived mission; otherwise scans roadmaps for one that contains it.
func findRoadmapForArchivedMission(spaceDir, archivedID string) (rid string, r map[string]interface{}, ok bool) {
	roadmapsDir := filepath.Join(spaceDir, "roadmaps")
	if data, err := os.ReadFile(filepath.Join(spaceDir, "current-roadmap")); err == nil {
		cur := strings.TrimSpace(string(data))
		if cur != "" {
			if rm, err := loadRoadmap(roadmapsDir, cur); err == nil && roadmapContainsMission(rm, archivedID) {
				return cur, rm, true
			}
		}
	}
	entries, err := os.ReadDir(roadmapsDir)
	if err != nil {
		return "", nil, false
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		ids = append(ids, strings.TrimSuffix(e.Name(), ".json"))
	}
	sort.Strings(ids)
	for _, id := range ids {
		rm, err := loadRoadmap(roadmapsDir, id)
		if err != nil {
			continue
		}
		if roadmapContainsMission(rm, archivedID) {
			return id, rm, true
		}
	}
	return "", nil, false
}

// suggestNextAfterArchive clears stale current, advances to the next incomplete
// roadmap mission when one exists, and prints a human hint.
func suggestNextAfterArchive(spaceDir, archivedID string) {
	if readCurrent(spaceDir) == archivedID {
		clearCurrent(spaceDir)
	}
	rid, r, found := findRoadmapForArchivedMission(spaceDir, archivedID)
	if !found {
		return
	}
	nextID, desc, state, ok := nextIncompleteOnRoadmap(spaceDir, r)
	if !ok {
		return
	}
	_ = writeCurrent(spaceDir, nextID)
	fmt.Printf("Next mission on roadmap %s: %s: %s (state=%s)\n", rid, nextID, desc, state)
	fmt.Printf("Suggested: new session → /sc-discuss %s (then /sc-run)\n", nextID)
}

// missionIncompleteState returns the display state for an incomplete mission, or "" if complete.
// Missing mission dir / mission.json → "pending". ready/shipped → complete ("").
func missionIncompleteState(missionPath string) string {
	data, err := os.ReadFile(filepath.Join(missionPath, "mission.json"))
	if err != nil {
		return "pending"
	}
	var mj map[string]interface{}
	if err := json.Unmarshal(data, &mj); err != nil {
		return "pending"
	}
	state, _ := mj["state"].(string)
	if state == "" {
		return "pending"
	}
	if state == "ready" || state == "shipped" {
		return ""
	}
	return state
}

func mapUse(args []string, spaceDir, roadmapsDir string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft map use <roadmap-id>")
		return 1
	}
	id := args[0]
	if _, err := loadRoadmap(roadmapsDir, id); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}
	if err := os.MkdirAll(spaceDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(spaceDir, "current-roadmap"), []byte(id+"\n"), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map:", err)
		return 1
	}
	fmt.Printf("Selected roadmap %s\n", id)
	return 0
}

func mapCurrent(spaceDir string) int {
	data, err := os.ReadFile(filepath.Join(spaceDir, "current-roadmap"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft map: no current roadmap")
		return 1
	}
	id := strings.TrimSpace(string(data))
	if id == "" {
		fmt.Fprintln(os.Stderr, "spacecraft map: no current roadmap")
		return 1
	}
	fmt.Println(id)
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
