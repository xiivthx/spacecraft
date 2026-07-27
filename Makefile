.PHONY: build install install-cli install-project install-global smoke uninstall clean help \
        test test-go test-config test-install test-gen-user-rules test-judge-break gate

ROOT      := $(CURDIR)
BIN       := $(ROOT)/spacecraft
LOCAL_BIN := $(HOME)/.local/bin
GLOBAL    := $(HOME)/.cursor
PROJECT   ?= .

help:
	@echo "Spacecraft targets:"
	@echo "  build           Build cmd/spacecraft -> ./spacecraft"
	@echo "  test            Full verification: go tests + vet + gofmt + config + judge-break + install smoke"
	@echo "  test-go         Go unit tests, go vet, gofmt check (cmd/spacecraft)"
	@echo "  test-config     Static config smoke (mcp/hooks JSON, frontmatter, no commands/)"
	@echo "  test-judge-break  Known-bad closeout fixtures must be rejected"
	@echo "  test-install    Bootstrap/install smoke into a throwaway temp dir"
	@echo "  test-gen-user-rules  RED/GREEN test for scripts/gen-user-rules.sh"
	@echo "  gate            Go tests + Cursor hook unit tests (hooks_test.sh)"
	@echo "  install         build + link CLI into ~/.local/bin + smoke check"
	@echo "  install-project Install full .cursor surface into PROJECT=<dir> (default .)"
	@echo "  install-global  Install agents + skills + MCP into ~/.cursor and link CLI (careful)"
	@echo "  smoke           Run post-install smoke checks (PROJECT=<dir>)"
	@echo "  uninstall       Remove CLI link + spacecraft MCP/agents/skills from ~/.cursor"
	@echo "  clean           Remove built binary"

# Full verification suite. Runs everything CI runs; humans use `make test`.
test: test-go test-config test-judge-break test-install
	@echo "All tests passed."

# Local pre-ship / PR gate: Go unit tests plus Cursor hook assertions.
gate: test-go
	bash $(ROOT)/.cursor/hooks/hooks_test.sh

test-go:
	go -C $(ROOT)/cmd/spacecraft test ./...
	go -C $(ROOT)/cmd/spacecraft vet ./...
	@unformatted=$$(gofmt -l $(ROOT)/cmd/spacecraft); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needs to run on:"; echo "$$unformatted"; exit 1; \
	fi
	@echo "gofmt clean"

test-config:
	@sh $(ROOT)/scripts/config-smoke.sh "$(ROOT)"

# Deterministic judge-break: known-bad fixtures must fail closeout-check.
test-judge-break: build
	@sh $(ROOT)/scripts/check-judge-break.sh "$(ROOT)" "$(BIN)"

# Install the surface + build the CLI into a throwaway dir, then smoke it.
# Temp dir lives under the repo (.tmp/, gitignored) so it works both in a
# sandboxed local run and in CI where $TMPDIR may be off-limits.
test-install: build
	@sh $(ROOT)/scripts/test-install.sh "$(ROOT)" "$(BIN)"

# Focused RED/GREEN test for the USER-RULES generator (no build required).
test-gen-user-rules:
	@sh $(ROOT)/scripts/test-gen-user-rules.sh "$(ROOT)"

build:
	cd $(ROOT)/cmd/spacecraft && go build -o $(BIN) .

# Safe default: build the CLI, link it, and validate. Config stays project-local
# (this repo already carries .cursor/ and .space/).
install: build install-cli smoke

install-cli: build
	@mkdir -p $(LOCAL_BIN)
	@ln -sf $(BIN) $(LOCAL_BIN)/spacecraft
	@echo "cli -> $(LOCAL_BIN)/spacecraft"

# Copy the full .cursor surface + .space scaffold into another project.
install-project:
	@sh $(ROOT)/scripts/install-cursor.sh "$(PROJECT)" "$(ROOT)"
	@sh $(ROOT)/scripts/smoke.sh "$(PROJECT)" "$(BIN)"

# Global Cursor install. Never writes ~/.cursorrules. Copies agents + skills,
# merges MCP servers; leaves unrelated global config untouched.
install-global: build install-cli
	@sh $(ROOT)/scripts/global-sync.sh "$(ROOT)" "$(GLOBAL)" agents install
	@echo "agents -> $(GLOBAL)/agents"
	@sh $(ROOT)/scripts/global-sync.sh "$(ROOT)" "$(GLOBAL)" skills install
	@echo "skills -> $(GLOBAL)/skills (sc-*)"
	@test -f "$(GLOBAL)/skills/sc-run/SKILL.md"
	@test -f "$(GLOBAL)/skills/sc-ship/SKILL.md"
	@test -f "$(GLOBAL)/skills/sc-quick/SKILL.md"
	@test -f "$(GLOBAL)/skills/sc-storm/SKILL.md"
	@test -f "$(GLOBAL)/skills/sc-discuss/references/lens-pass.md"
	@test -f "$(GLOBAL)/skills/sc-judge/references/judge-break/open-issues/expect.json"
	@python3 $(ROOT)/scripts/mcp-merge.py merge "$(GLOBAL)/mcp.json" "$(ROOT)/.cursor/mcp.json"
	@sh $(ROOT)/scripts/install-global-hooks.sh "$(ROOT)" "$(GLOBAL)"
	@sh $(ROOT)/scripts/gen-user-rules.sh "$(ROOT)/.cursor/rules" "$(GLOBAL)/spacecraft/USER-RULES.txt"
	@echo "user rules -> $(GLOBAL)/spacecraft/USER-RULES.txt (paste into Settings > Rules > User Rules)"
	@echo "Global install complete. Restart Cursor to pick up /sc-run, /sc-ship, /sc-quick, and sc-storm."

smoke:
	@sh $(ROOT)/scripts/smoke.sh "$(PROJECT)" "$(BIN)"

# Remove only what we added; never clobber unrelated user MCP servers or config.
uninstall:
	@if [ -L "$(LOCAL_BIN)/spacecraft" ]; then rm -f "$(LOCAL_BIN)/spacecraft" && echo "removed $(LOCAL_BIN)/spacecraft"; fi
	@sh $(ROOT)/scripts/global-sync.sh "$(ROOT)" "$(GLOBAL)" agents uninstall
	@echo "removed spacecraft agents from $(GLOBAL)/agents"
	@sh $(ROOT)/scripts/global-sync.sh "$(ROOT)" "$(GLOBAL)" skills uninstall
	@echo "removed spacecraft skills from $(GLOBAL)/skills"
	@if [ -f "$(GLOBAL)/mcp.json" ]; then python3 $(ROOT)/scripts/mcp-merge.py unmerge "$(GLOBAL)/mcp.json" "$(ROOT)/.cursor/mcp.json"; fi
	@echo "Uninstall complete. Project-local .cursor/ and .space/ left untouched."

clean:
	rm -f $(BIN)
