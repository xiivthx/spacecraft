.PHONY: help test build sc-design-serve sc-design-open install uninstall reinstall

OPENACODE_GLOBAL ?= $(HOME)/.config/opencode

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ── global config symlink targets ─────────────────────────────────────────

install: ## Symlink spacecraft skills/agents/commands + docs to ~/.config/opencode (merges agents into opencode.json)
	@if [ -d "$(CURDIR)/.opencode/skills" ] && [ "$$FORCE" != "1" ]; then \
		echo "WARNING: this project already has .opencode/ locally."; \
		echo "  Global install may cause double-loading in this project."; \
		echo "  Use FORCE=1 to override, or run from another project."; \
		exit 1; \
	fi
	@echo "=== Installing spacecraft to $(OPENACODE_GLOBAL) ==="
	@test -f opencode.json || (echo "ERROR: run from spacecraft repo root" && exit 1)
	@for dir in skills agents commands; do \
		if [ -e "$(OPENACODE_GLOBAL)/$$dir" ] && [ ! -L "$(OPENACODE_GLOBAL)/$$dir" ]; then \
			echo "  Backing up $$dir/ -> $$dir.bak"; \
			mv "$(OPENACODE_GLOBAL)/$$dir" "$(OPENACODE_GLOBAL)/$$dir.bak"; \
		fi; \
		ln -sfn "$(CURDIR)/.opencode/$$dir" "$(OPENACODE_GLOBAL)/$$dir"; \
		echo "  $$dir/ -> $(CURDIR)/.opencode/$$dir"; \
	done
	@if [ -e "$(OPENACODE_GLOBAL)/AGENTS.md" ] && [ ! -L "$(OPENACODE_GLOBAL)/AGENTS.md" ]; then \
		echo "  Backing up AGENTS.md -> AGENTS.md.bak"; \
		mv "$(OPENACODE_GLOBAL)/AGENTS.md" "$(OPENACODE_GLOBAL)/AGENTS.md.bak"; \
	fi; \
	ln -sf "$(CURDIR)/.opencode/AGENTS.md" "$(OPENACODE_GLOBAL)/AGENTS.md"; \
	echo "  AGENTS.md -> $(CURDIR)/.opencode/AGENTS.md"
	@if [ -L "$(OPENACODE_GLOBAL)/opencode.json" ]; then \
		echo "  Global opencode.json is symlinked — skipping merge"; \
	elif jq -e '.agent' "$(OPENACODE_GLOBAL)/opencode.json" >/dev/null 2>&1; then \
		echo "  Agent definitions already in global opencode.json — skipping merge"; \
	elif [ -f "$(OPENACODE_GLOBAL)/opencode.json" ]; then \
		echo "  Merging agent definitions into global opencode.json..."; \
		[ ! -f "$(OPENACODE_GLOBAL)/opencode.json.bak" ] && cp "$(OPENACODE_GLOBAL)/opencode.json" "$(OPENACODE_GLOBAL)/opencode.json.bak"; \
		jq -s '.[0] * {agent: .[1].agent}' "$(OPENACODE_GLOBAL)/opencode.json" "$(CURDIR)/opencode.json" > "$(OPENACODE_GLOBAL)/opencode.json.tmp" && \
		mv "$(OPENACODE_GLOBAL)/opencode.json.tmp" "$(OPENACODE_GLOBAL)/opencode.json"; \
		echo "  opencode.json merged (backup at opencode.json.bak)"; \
	else \
		echo "  No global opencode.json found — copying project config"; \
		cp "$(CURDIR)/opencode.json" "$(OPENACODE_GLOBAL)/opencode.json"; \
	fi
	@echo "=== Done. Restart OpenCode for changes to take effect. ==="

uninstall: ## Remove spacecraft symlinks from ~/.config/opencode
	@echo "=== Uninstalling spacecraft from $(OPENACODE_GLOBAL) ==="
	@for dir in skills agents commands; do \
		if [ -L "$(OPENACODE_GLOBAL)/$$dir" ]; then \
			rm "$(OPENACODE_GLOBAL)/$$dir"; \
			echo "  Removed $$dir symlink"; \
			if [ -d "$(OPENACODE_GLOBAL)/$$dir.bak" ]; then \
				mv "$(OPENACODE_GLOBAL)/$$dir.bak" "$(OPENACODE_GLOBAL)/$$dir"; \
				echo "  Restored $$dir from backup"; \
			fi; \
		fi; \
	done
	@if [ -L "$(OPENACODE_GLOBAL)/AGENTS.md" ]; then \
		rm "$(OPENACODE_GLOBAL)/AGENTS.md"; \
		echo "  Removed AGENTS.md symlink"; \
		if [ -f "$(OPENACODE_GLOBAL)/AGENTS.md.bak" ]; then \
			mv "$(OPENACODE_GLOBAL)/AGENTS.md.bak" "$(OPENACODE_GLOBAL)/AGENTS.md"; \
			echo "  Restored AGENTS.md from backup"; \
		fi; \
	fi
	@if [ -f "$(OPENACODE_GLOBAL)/opencode.json.bak" ]; then \
		echo "  opencode.json.bak exists (original backup) — restore manually: mv ~/.config/opencode/opencode.json.bak ~/.config/opencode/opencode.json"; \
	fi
	@echo "=== Done. ==="

reinstall: uninstall install ## Uninstall then reinstall

fix-opencode: ## Strip project-specific keys from global opencode.json (if previously merged with full project config)
	@if [ -f "$(OPENACODE_GLOBAL)/opencode.json" ] && [ ! -L "$(OPENACODE_GLOBAL)/opencode.json" ]; then \
		jq 'del(.default_agent, .share, .lsp, .instructions)' "$(OPENACODE_GLOBAL)/opencode.json" > "$(OPENACODE_GLOBAL)/opencode.json.tmp" && \
		mv "$(OPENACODE_GLOBAL)/opencode.json.tmp" "$(OPENACODE_GLOBAL)/opencode.json" && \
		echo "  Stripped default_agent/share/lsp/instructions from global opencode.json"; \
	else \
		echo "  Skipped (symlinked or missing)"; \
	fi

clean-global: ## Remove all spacecraft symlinks + restore opencode.json.bak (for resetting global config)
	@echo "=== Cleaning spacecraft from $(OPENACODE_GLOBAL) ==="
	@for dir in skills agents commands; do \
		if [ -L "$(OPENACODE_GLOBAL)/$$dir" ]; then \
			rm "$(OPENACODE_GLOBAL)/$$dir"; \
			echo "  Removed $$dir symlink"; \
		fi; \
	done
	@for file in AGENTS.md PERSONA.md; do \
		if [ -L "$(OPENACODE_GLOBAL)/$$file" ]; then \
			rm "$(OPENACODE_GLOBAL)/$$file"; \
			echo "  Removed $$file symlink"; \
		fi; \
	done
	@if [ -f "$(OPENACODE_GLOBAL)/opencode.json.bak" ]; then \
		mv "$(OPENACODE_GLOBAL)/opencode.json.bak" "$(OPENACODE_GLOBAL)/opencode.json"; \
		echo "  Restored opencode.json from backup"; \
	fi
	@echo "=== Clean. Run 'make install' to re-symlink. ==="

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
