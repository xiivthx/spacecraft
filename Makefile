.PHONY: build install install-cli install-project install-global smoke uninstall clean help \
        test test-go test-config test-install

ROOT      := $(CURDIR)
BIN       := $(ROOT)/spacecraft
LOCAL_BIN := $(HOME)/.local/bin
GLOBAL    := $(HOME)/.cursor
PROJECT   ?= .

help:
	@echo "Spacecraft targets:"
	@echo "  build           Build cmd/spacecraft -> ./spacecraft"
	@echo "  test            Full verification: go tests + vet + gofmt + config + install smoke"
	@echo "  test-go         Go unit tests, go vet, gofmt check (cmd/spacecraft)"
	@echo "  test-config     Static config smoke (mcp/hooks JSON, frontmatter, no commands/)"
	@echo "  test-install    Bootstrap/install smoke into a throwaway temp dir"
	@echo "  install         build + link CLI into ~/.local/bin + smoke check"
	@echo "  install-project Install full .cursor surface into PROJECT=<dir> (default .)"
	@echo "  install-global  Merge agents + MCP into ~/.cursor and link CLI (careful)"
	@echo "  smoke           Run post-install smoke checks (PROJECT=<dir>)"
	@echo "  uninstall       Remove CLI link + spacecraft MCP servers + global agents"
	@echo "  clean           Remove built binary"

# Full verification suite. Runs everything CI runs; humans use `make test`.
test: test-go test-config test-install
	@echo "All tests passed."

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

# Install the surface + build the CLI into a throwaway dir, then smoke it.
# Temp dir lives under the repo (.tmp/, gitignored) so it works both in a
# sandboxed local run and in CI where $TMPDIR may be off-limits.
# Seeds a pre-existing unrelated hook so merge must preserve user hooks
# alongside the spacecraft ship-gate hook. Never writes ~/.cursorrules.
test-install: build
	@mkdir -p $(ROOT)/.tmp; \
	tmp=$$(mktemp -d "$(ROOT)/.tmp/install-smoke.XXXXXX"); \
	fake_home=$$tmp/home; \
	trap 'rm -rf "$$tmp"' EXIT INT TERM; \
	echo "Install smoke in $$tmp"; \
	mkdir -p "$$tmp/.cursor" "$$fake_home"; \
	printf '%s\n' '{"version":1,"hooks":{"beforeShellExecution":[{"command":".cursor/hooks/user-unrelated.sh","matcher":"echo"}]}}' \
		> "$$tmp/.cursor/hooks.json"; \
	HOME="$$fake_home" sh $(ROOT)/scripts/install-cursor.sh "$$tmp" "$(ROOT)"; \
	sh $(ROOT)/scripts/smoke.sh "$$tmp" "$(BIN)"; \
	hooks="$$tmp/.cursor/hooks.json"; \
	if ! grep -q 'user-unrelated.sh' "$$hooks"; then \
		echo "FAIL: install clobbered pre-existing hooks (user-unrelated.sh missing)"; \
		exit 1; \
	fi; \
	if ! grep -q 'check-ship-commands.sh' "$$hooks"; then \
		echo "FAIL: install missing spacecraft ship-gate hook (check-ship-commands.sh)"; \
		exit 1; \
	fi; \
	echo "  ok   hooks merge preserves user + ship-gate"; \
	if [ -f "$$fake_home/.cursorrules" ]; then \
		echo "FAIL: install wrote legacy ~/.cursorrules"; \
		exit 1; \
	fi; \
	echo "  ok   no legacy ~/.cursorrules"

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

# Global Cursor install. Never writes ~/.cursorrules. Merges MCP servers and
# copies agents; leaves unrelated global config untouched.
install-global: build install-cli
	@mkdir -p $(GLOBAL)/agents
	@cp $(ROOT)/.cursor/agents/sc-*.md $(GLOBAL)/agents/
	@echo "agents -> $(GLOBAL)/agents"
	@python3 $(ROOT)/scripts/mcp-merge.py merge $(GLOBAL)/mcp.json $(ROOT)/.cursor/mcp.json
	@echo "Global install complete. Restart Cursor to pick up config."

smoke:
	@sh $(ROOT)/scripts/smoke.sh "$(PROJECT)" "$(BIN)"

# Remove only what we added; never clobber unrelated user MCP servers or config.
uninstall:
	@if [ -L $(LOCAL_BIN)/spacecraft ]; then rm -f $(LOCAL_BIN)/spacecraft && echo "removed $(LOCAL_BIN)/spacecraft"; fi
	@for a in $(ROOT)/.cursor/agents/sc-*.md; do rm -f $(GLOBAL)/agents/$$(basename $$a); done
	@echo "removed spacecraft agents from $(GLOBAL)/agents"
	@if [ -f $(GLOBAL)/mcp.json ]; then python3 $(ROOT)/scripts/mcp-merge.py unmerge $(GLOBAL)/mcp.json $(ROOT)/.cursor/mcp.json; fi
	@echo "Uninstall complete. Project-local .cursor/ and .space/ left untouched."

clean:
	rm -f $(BIN)
