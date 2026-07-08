.PHONY: help test build sc-design-serve sc-design-open

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: | build ## Run Go tests and Node integration tests
	cd scripts/src && go test ./...
	node --test tests/*.test.mjs

build: ## Build the Go binary
	cd scripts/src && go build -o ../spacecraft .

sc-design-serve: ## Serve design HTML artifacts (use ARGS="[artifact-or-dir]")
	node .opencode/skills/sc-design/scripts/serve-html.mjs $(ARGS)

sc-design-open: ## Serve and open design HTML artifacts in browser (use ARGS="[artifact-or-dir]")
	node .opencode/skills/sc-design/scripts/serve-html.mjs --open $(ARGS)

status: build ## Show active mission status
	./scripts/spacecraft status

resume: build ## Resume active mission (CLI snapshot — use /sc-resume for full handoff)
	@echo "═══ Mission status ═══"
	@./scripts/spacecraft status || true
	@echo ""
	@echo "═══ Workflow ═══"
	@./scripts/spacecraft flow || true
	@echo ""
	@echo "═══ Git ═══"
	@./scripts/spacecraft git-info || true
	@echo ""
	@echo "═══ Last commit ═══"
	@git log -1 --oneline --no-decorate 2>/dev/null || echo "(no commits)"
	@echo ""
	@echo "Tip: use /sc-resume in the TUI for full handoff resume with next-action guidance."
