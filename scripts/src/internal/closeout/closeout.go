// Package closeout provides release readiness checks for mission closeout.
// All functions return errors instead of calling os.Exit.
package closeout

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"spacecraft/internal/mission"
)

// defaultReleaseGateStatuses lists statuses that satisfy a release gate requirement.
var defaultReleaseGateStatuses = map[string]bool{
	"bumped": true, "checked": true, "complete": true, "completed": true,
	"deferred": true, "done": true, "passed": true, "present": true, "updated": true,
}

// tagPlanReleaseGateStatuses extends default statuses with "planned".
var tagPlanReleaseGateStatuses = map[string]bool{
	"bumped": true, "checked": true, "complete": true, "completed": true,
	"deferred": true, "done": true, "passed": true, "present": true, "updated": true,
	"planned": true,
}

// Runner provides git command execution for closeout checks.
type Runner interface {
	Run(name string, args ...string) (exitCode int, stdout, stderr string)
}

// Checker performs closeout readiness checks.
type Checker struct {
	store     mission.MissionStore
	gitRunner Runner
}

// NewChecker creates a new Checker.
func NewChecker(store mission.MissionStore, gitRunner Runner) *Checker {
	return &Checker{store: store, gitRunner: gitRunner}
}

// CheckResult holds the results of a closeout check.
type CheckResult struct {
	Errors   []string
	Warnings []string
}

// Check performs all closeout readiness checks for a mission.
func (c *Checker) Check(id string, m *mission.Mission, plan *mission.Plan, review *mission.Review, evidenceCount int) CheckResult {
	var errors []string
	var warnings []string

	// Clarification check
	if m.Clarification.Status == "open" || m.Clarification.BlockingQuestions > 0 {
		errors = append(errors, "Resolve blocking clarification questions.")
	}

	// Task completion check
	var incomplete []mission.Task
	if plan != nil {
		for _, t := range plan.Tasks {
			if t.Status == nil || *t.Status != "completed" {
				incomplete = append(incomplete, t)
			}
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

	// Evidence check
	if evidenceCount == 0 {
		errors = append(errors, "Capture verification evidence in evidence.jsonl.")
	}

	// Review checks
	if review != nil {
		if review.Status == nil || *review.Status != "ready" {
			stat := "missing"
			if review.Status != nil {
				stat = *review.Status
			}
			errors = append(errors, fmt.Sprintf("Review status must be ready; current status is %s.", stat))
		}

		blocking := blockingReviewFindings(review)
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

	// Git checks
	gitCode, gitOut, _ := c.gitRunner.Run("git", "rev-parse", "--is-inside-work-tree")
	inRepo := gitCode == 0 && strings.TrimSpace(gitOut) == "true"

	if !inRepo {
		errors = append(errors, "Run closeout inside a git worktree.")
	} else {
		_, branchOut, _ := c.gitRunner.Run("git", "branch", "--show-current")
		branch := strings.TrimSpace(branchOut)

		if branch == "" || branch == "main" {
			errors = append(errors, "Closeout must run from a non-main work branch.")
		}

		// Check dirty
		statusCode, statusOut, _ := c.gitRunner.Run("git", "status", "--short")
		statusText := strings.TrimSpace(statusOut)
		var dirtyFiles int
		if statusText != "" {
			dirtyFiles = len(strings.Split(statusText, "\n"))
		}
		if statusCode == 0 && dirtyFiles > 0 {
			errors = append(errors, fmt.Sprintf("Commit, stash, or remove dirty worktree changes (%d files).", dirtyFiles))
		}

		// Check rebase status
		mainAncestorCode, _, _ := c.gitRunner.Run("git", "merge-base", "--is-ancestor", "main", "HEAD")
		if mainAncestorCode != 0 {
			errors = append(errors, "Rebase the work branch on latest main, then rerun verification.")
		}

		// Check commit count
		commitCountCode, commitCountOut, _ := c.gitRunner.Run("git", "rev-list", "--count", "main..HEAD")
		count := -1
		if commitCountCode == 0 {
			c, err := strconv.Atoi(strings.TrimSpace(commitCountOut))
			if err == nil {
				count = c
				if count > 5 {
					errors = append(errors, fmt.Sprintf("Squash/fixup branch history to 5 or fewer final commits; current count is %d.", count))
				}
			}
		} else {
			warnings = append(warnings, "Could not count commits from main..HEAD.")
		}

		// Check commit subjects
		commitSubjectsCode, commitSubjectsOut, _ := c.gitRunner.Run("git", "log", "--format=%s", "main..HEAD")
		if commitSubjectsCode == 0 {
			out := strings.TrimSpace(commitSubjectsOut)
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

	return CheckResult{Errors: errors, Warnings: warnings}
}

// ReleaseGateSatisfied checks if a release gate has an allowed status.
func ReleaseGateSatisfied(gate *mission.ReleaseGate, allowedStatuses map[string]bool) bool {
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

func releaseReadinessErrors(rr mission.ReleaseReadiness) []string {
	var errors []string
	if !ReleaseGateSatisfied(rr.Version, defaultReleaseGateStatuses) {
		errors = append(errors, "Record version bump or explicit deferral with rationale in review.json releaseReadiness.version.")
	}
	if !ReleaseGateSatisfied(rr.Changelog, defaultReleaseGateStatuses) {
		errors = append(errors, "Record changelog update or explicit deferral with rationale in review.json releaseReadiness.changelog.")
	}
	if !ReleaseGateSatisfied(rr.SpecNote, defaultReleaseGateStatuses) {
		errors = append(errors, "Record short spec/release note update or explicit deferral with rationale in review.json releaseReadiness.specNote.")
	}
	if !ReleaseGateSatisfied(rr.TagPlan, tagPlanReleaseGateStatuses) {
		errors = append(errors, "Record the post-merge version tag plan in review.json releaseReadiness.tagPlan.")
	}
	if !ReleaseGateSatisfied(rr.PostRebaseVerification, defaultReleaseGateStatuses) {
		errors = append(errors, "Record verification after latest rebase in review.json releaseReadiness.postRebaseVerification.")
	}
	return errors
}

func blockingReviewFindings(review *mission.Review) []mission.Finding {
	if review == nil {
		return nil
	}
	var blocking []mission.Finding
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
