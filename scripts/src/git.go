package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"spacecraft/internal/gitutil"
)

func gitInfo() GitInfoData {
	return gitutil.GitInfo(gitutil.OSCommandRunner{})
}

func printGitInfo() {
	git := gitInfo()
	if !git.Available {
		fmt.Println("Git: command unavailable")
		return
	}
	if !git.IsRepo {
		fmt.Println("Git: not a git worktree")
		return
	}

	fmt.Println("Git: worktree detected")
	root := "(unknown)"
	if git.Root != "" {
		root = git.Root
	}
	fmt.Printf("Root: %s\n", root)
	branch := "(detached)"
	if git.Branch != "" {
		branch = git.Branch
	}
	fmt.Printf("Branch: %s\n", branch)
	sha := "(no commit)"
	if git.Sha != "" {
		sha = git.Sha
	}
	fmt.Printf("HEAD: %s\n", sha)
	status := "clean"
	if git.Dirty {
		status = fmt.Sprintf("dirty (%d files)", git.DirtyFiles)
	}
	fmt.Printf("Status: %s\n", status)
}

func printGitSuggestion(args []string) {
	resolution := resolveMission("")
	if resolution.Safety != "safe" || resolution.Selected == nil {
		fail(formatResolutionBlock(resolution, "git-suggest"))
	}

	id := resolution.Selected.ID
	var mission *Mission
	err := readJson(filepath.Join(missionDir(id), "mission.json"), &mission)
	missionPart := strings.ToLower(id)

	branchTypes := map[string]bool{
		"feat": true, "fix": true, "docs": true, "refactor": true,
		"test": true, "build": true, "ci": true, "chore": true,
		"perf": true, "style": true, "issue": true, "release": true,
	}

	requestedType := ""
	if len(args) > 0 {
		requestedType = strings.ToLower(args[0])
	}

	typ := "feat"
	var slugParts []string
	if branchTypes[requestedType] {
		typ = requestedType
		slugParts = args[1:]
	} else {
		slugParts = args
	}

	slugBase := strings.Join(slugParts, " ")
	if slugBase == "" {
		if err == nil && mission != nil && mission.Title != "" {
			slugBase = mission.Title
		} else {
			slugBase = "mission"
		}
	}
	slug := slugify(slugBase)
	// Strip mission id prefix from slug to avoid duplication:
	// e.g. "m07fsgjf6-makefile-migration" becomes "makefile-migration"
	missionPrefix := strings.ToLower(id)
	if strings.HasPrefix(slug, missionPrefix+"-") {
		slug = slug[len(missionPrefix)+1:]
	} else if strings.HasPrefix(slug, missionPrefix) && slug != missionPrefix {
		slug = slug[len(missionPrefix):]
	}
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "mission"
	}

	versionStr := strings.TrimSpace(strings.Join(slugParts, " "))
	versionStr = strings.ToLower(versionStr)
	if !strings.HasPrefix(versionStr, "v") {
		versionStr = "v" + versionStr
	}
	versionStr = regexpReplace(`[^a-z0-9._-]+`, "-", versionStr)
	versionStr = strings.Trim(versionStr, "-")
	if versionStr == "" || versionStr == "v" {
		versionStr = "v0.1.0"
	}

	branch := fmt.Sprintf("%s/%s/%s", typ, missionPart, slug)
	if typ == "release" {
		branch = fmt.Sprintf("release/%s", versionStr)
	}

	commitType := typ
	if typ == "issue" || typ == "release" {
		commitType = "chore"
	}

	fmt.Println("Spacecraft git strategy: release branching")
	if resolution.Selected != nil {
		fmt.Printf("Mission: %s (%s)\n", resolution.Selected.Title, resolution.Selected.ID)
		src := "unknown"
		if resolution.Source != nil {
			src = *resolution.Source
		}
		fmt.Printf("Selected by: %s\n", src)
	}
	fmt.Println("Base: latest main")
	fmt.Println("Main: protected; no direct writes")
	fmt.Printf("Branch: %s\n", branch)
	if typ == "release" {
		fmt.Println("Rule: release branches are for version/changelog/spec preparation only.")
	} else {
		fmt.Println("Rule: one branch per feature/issue/tightly scoped change.")
	}
	fmt.Println("Final branch history: target 1-3 commits; maximum 5 unless justified.")
	fmt.Println("Before merge: rebase on latest main, verify, bump version, update changelog/spec.")
	fmt.Println("Merge: git merge --no-ff <branch>")
	fmt.Println("After merge: create annotated version tag.")
	fmt.Println("Use a separate worktree for large, risky, or multi-session branches.")
	fmt.Println("")
	fmt.Println("Conventional Commits format:")
	fmt.Println("<type>[optional scope]: <description>")
	fmt.Println("")
	fmt.Println("Common types: feat, fix, docs, refactor, test, build, ci, chore, perf, style, revert")
	fmt.Println("Examples:")
	if typ == "release" {
		fmt.Printf("chore: prepare release %s\n", versionStr)
	} else {
		fmt.Printf("%s: add focused mission change\n", commitType)
	}
	fmt.Println("docs: update changelog for v0.2.0")
	fmt.Println("chore: bump version to v0.2.0")
}
