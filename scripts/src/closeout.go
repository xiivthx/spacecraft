package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var defaultReleaseGateStatuses = map[string]bool{
	"bumped": true, "checked": true, "complete": true, "completed": true,
	"deferred": true, "done": true, "passed": true, "present": true, "updated": true,
}

var tagPlanReleaseGateStatuses = map[string]bool{
	"bumped": true, "checked": true, "complete": true, "completed": true,
	"deferred": true, "done": true, "passed": true, "present": true, "updated": true,
	"planned": true,
}

func releaseGateSatisfied(gate *ReleaseGate, allowedStatuses map[string]bool) bool {
	if gate == nil || gate.Status == nil {
		return false
	}
	status := strings.ToLower(strings.TrimSpace(*gate.Status))
	if !allowedStatuses[status] {
		return false
	}
	if status == "deferred" {
		if gate.Rationale == nil || strings.TrimSpace(*gate.Rationale) == "" {
			return false
		}
	}
	return true
}

func releaseReadinessErrors(rr ReleaseReadiness) []string {
	var errors []string
	if !releaseGateSatisfied(rr.Version, defaultReleaseGateStatuses) {
		errors = append(errors, "Record version bump or explicit deferral with rationale in review.json releaseReadiness.version.")
	}
	if !releaseGateSatisfied(rr.Changelog, defaultReleaseGateStatuses) {
		errors = append(errors, "Record changelog update or explicit deferral with rationale in review.json releaseReadiness.changelog.")
	}
	if !releaseGateSatisfied(rr.SpecNote, defaultReleaseGateStatuses) {
		errors = append(errors, "Record short spec/release note update or explicit deferral with rationale in review.json releaseReadiness.specNote.")
	}
	if !releaseGateSatisfied(rr.TagPlan, tagPlanReleaseGateStatuses) {
		errors = append(errors, "Record the post-merge version tag plan in review.json releaseReadiness.tagPlan.")
	}
	if !releaseGateSatisfied(rr.PostRebaseVerification, defaultReleaseGateStatuses) {
		errors = append(errors, "Record verification after latest rebase in review.json releaseReadiness.postRebaseVerification.")
	}
	return errors
}

func blockingReviewFindings(review *Review) []Finding {
	if review == nil {
		return nil
	}
	var blocking []Finding
	for _, f := range review.Findings {
		blocks := f.BlocksShip != nil && *f.BlocksShip
		critical := f.Severity != nil && *f.Severity == "critical"
		if blocks || critical {
			blocking = append(blocking, f)
		}
	}
	return blocking
}

func conventionalCommitSubject(subject string) bool {
	re := regexp.MustCompile(`^(feat|fix|docs|refactor|test|build|ci|chore|perf|style|revert)(\([a-z0-9._/-]+\))?!?: .+`)
	return re.MatchString(subject)
}


func releaseCloseoutCheck() {
	var errors []string
	var warnings []string
	res := requireResolvedMission("closeout-check")
	id := res.Selected.ID
	dir := missionDir(id)

	var mission Mission
	errM := readJson(filepath.Join(dir, "mission.json"), &mission)

	var plan Plan
	errP := readJson(filepath.Join(dir, "plan.json"), &plan)

	var review Review
	errR := readJson(filepath.Join(dir, "review.json"), &review)

	evCount := countEvidence(filepath.Join(dir, "evidence.jsonl"))

	if errM != nil {
		errors = append(errors, "Missing mission.json")
	}
	if errP != nil {
		errors = append(errors, "Missing plan.json")
	}
	if errR != nil {
		errors = append(errors, "Missing review.json")
	}

	if mission.Clarification.Status == "open" || mission.Clarification.BlockingQuestions > 0 {
		errors = append(errors, "Resolve blocking clarification questions.")
	}

	var incomplete []Task
	for _, t := range plan.Tasks {
		if t.Status == nil || *t.Status != "completed" {
			incomplete = append(incomplete, t)
		}
	}
	if len(incomplete) > 0 {
		names := []string{}
		for _, t := range incomplete {
			if t.ID != nil {
				names = append(names, *t.ID)
			} else if t.Title != nil {
				names = append(names, *t.Title)
			}
		}
		errors = append(errors, fmt.Sprintf("Complete plan tasks: %s.", strings.Join(names, ", ")))
	}

	if evCount == 0 {
		errors = append(errors, "Capture verification evidence in evidence.jsonl.")
	}

	if errR == nil {
		if review.Status == nil || *review.Status != "ready" {
			stat := "missing"
			if review.Status != nil {
				stat = *review.Status
			}
			errors = append(errors, fmt.Sprintf("Review status must be ready; current status is %s.", stat))
		}

		blocking := blockingReviewFindings(&review)
		if len(blocking) > 0 {
			names := []string{}
			for _, f := range blocking {
				if f.ID != nil {
					names = append(names, *f.ID)
				} else if f.Summary != nil {
					names = append(names, *f.Summary)
				} else {
					names = append(names, "unnamed")
				}
			}
			errors = append(errors, fmt.Sprintf("Fix blocking review findings: %s.", strings.Join(names, ", ")))
		}

		errors = append(errors, releaseReadinessErrors(review.ReleaseReadiness)...)
	}

	git := gitInfo()
	if !git.IsRepo {
		errors = append(errors, "Run closeout inside a git worktree.")
	} else {
		if git.Branch == "" || git.Branch == "main" {
			errors = append(errors, "Closeout must run from a non-main work branch.")
		}
		if git.Dirty {
			errors = append(errors, fmt.Sprintf("Commit, stash, or remove dirty worktree changes (%d files).", git.DirtyFiles))
		}

		mainAncestor := runCommand([]string{"git", "merge-base", "--is-ancestor", "main", "HEAD"})
		if mainAncestor.exitCode != 0 {
			errors = append(errors, "Rebase the work branch on latest main, then rerun verification.")
		}

		commitCount := runCommand([]string{"git", "rev-list", "--count", "main..HEAD"})
		count := -1
		if commitCount.exitCode == 0 {
			c, err := strconv.Atoi(strings.TrimSpace(commitCount.stdout))
			if err == nil {
				count = c
				if count > 5 {
					errors = append(errors, fmt.Sprintf("Squash/fixup branch history to 5 or fewer final commits; current count is %d.", count))
				}
			}
		} else {
			warnings = append(warnings, "Could not count commits from main..HEAD.")
		}

		commitSubjects := runCommand([]string{"git", "log", "--format=%s", "main..HEAD"})
		if commitSubjects.exitCode == 0 {
			out := strings.TrimSpace(commitSubjects.stdout)
			var subjects []string
			if out != "" {
				for _, line := range strings.Split(out, "\n") {
					if strings.TrimSpace(line) != "" {
						subjects = append(subjects, line)
					}
				}
			}
			if count == 0 || len(subjects) == 0 {
				errors = append(errors, "Create at least one final Conventional Commit before closeout.")
			}
			var invalid []string
			for _, s := range subjects {
				if !conventionalCommitSubject(s) {
					invalid = append(invalid, s)
				}
			}
			if len(invalid) > 0 {
				errors = append(errors, fmt.Sprintf("Fix non-Conventional Commit subjects: %s", strings.Join(invalid, "; ")))
			}
		} else {
			errors = append(errors, "Check final commit subjects before closeout.")
		}
	}

	if len(errors) > 0 {
		fmt.Printf("Spacecraft closeout blocked for %s:\n", id)
		for _, e := range errors {
			fmt.Printf("- %s\n", e)
		}
		if len(warnings) > 0 {
			fmt.Println("Warnings:")
			for _, w := range warnings {
				fmt.Printf("- %s\n", w)
			}
		}
		os.Exit(1)
	}

	fmt.Printf("Spacecraft closeout ready for %s.\n", id)
	fmt.Println("Next: rebase already satisfied, run final verification if needed, merge with git merge --no-ff <branch>, tag release, then delete merged branch unless kept.")
	if len(warnings) > 0 {
		fmt.Println("Warnings:")
		for _, w := range warnings {
			fmt.Printf("- %s\n", w)
		}
	}
}
