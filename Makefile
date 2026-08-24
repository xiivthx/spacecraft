.PHONY: build install install-cli install-project install-global install-machine smoke uninstall clean help \
        test test-cli test-config test-install test-gen-user-rules test-judge-break gate \
        sync-antigravity test-antigravity install-antigravity install-antigravity-global install-antigravity-project

ROOT      := $(CURDIR)
BIN       := $(ROOT)/spacecraft
LOCAL_BIN := $(HOME)/.local/bin
GLOBAL    := $(HOME)/.cursor
PROJECT   ?= .

help:
	@echo "Spacecraft targets:"
	@echo "  build           Link Node CLI (cli/spacecraft.mjs) -> ./spacecraft"
	@echo "  test            Full verification: Node CLI tests + config + judge-break + install smoke + antigravity"
	@echo "  test-cli        Node CLI unit tests (cli/test/)"
	@echo "  test-config     Static config smoke (mcp/hooks JSON, frontmatter, no commands/)"
	@echo "  test-antigravity Smoke check for Antigravity plugin, rules, skills, and hooks"
	@echo "  test-judge-break  Known-bad closeout fixtures must be rejected"
	@echo "  test-install    Bootstrap/install smoke into a throwaway temp dir"
	@echo "  test-gen-user-rules  RED/GREEN test for scripts/gen-user-rules.sh"
	@echo "  gate            Node CLI tests + Cursor hook unit tests (hooks_test.sh)"
	@echo "  install         build + link CLI into ~/.local/bin + smoke check"
	@echo "  install-project Install full .cursor surface into PROJECT=<dir> (default .)"
	@echo "  install-global  Install agents + lean-core skills + MCP into ~/.cursor (careful)"
	@echo "                  advanced escape hatch: SPACECRAFT_SKILL_PROFILE=full or FULL=1 (User encyclopedias; prefer project packs)"
	@echo "  install-machine Link Node CLI + one-shot User-layer install via scripts/install-machine.sh"
	@echo "  install-antigravity         Install global Spacecraft plugin for Antigravity"
	@echo "  install-antigravity-project Install Spacecraft into project .agents/ for Antigravity (PROJECT=<dir>)"
	@echo "  sync-antigravity            Regenerate/sync Antigravity plugin rules, skills, subagents, and hooks"
	@echo "  smoke           Run post-install smoke checks (PROJECT=<dir>)"
	@echo "  uninstall       Remove CLI link + spacecraft MCP/agents/skills from ~/.cursor"
	@echo "  clean           Remove ./spacecraft CLI link"

# Full verification suite. Runs everything CI runs; humans use `make test`.
test: test-cli test-config test-antigravity test-judge-break test-install
	@echo "All tests passed."

# Local pre-ship / PR gate: Node CLI tests plus Cursor hook assertions.
gate: test-cli
	bash $(ROOT)/.cursor/hooks/hooks_test.sh

test-cli:
	node --test $(ROOT)/cli/test/*.test.mjs

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
	chmod +x $(ROOT)/cli/spacecraft.mjs
	ln -sf $(ROOT)/cli/spacecraft.mjs $(BIN)

# Safe default: wire the Node CLI, link it, and validate. Config stays project-local
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

# Global Cursor install. Never writes ~/.cursorrules. Copies agents + skills
# (SPACECRAFT_SKILL_PROFILE=lean default). Prefer project packs for domain skills.
# Advanced escape hatch: FULL=1 or SPACECRAFT_SKILL_PROFILE=full installs User-layer
# domain encyclopedias (documented --full equivalent). Merges MCP servers; leaves
# unrelated global config untouched. Lean reconcile prunes spacecraft-managed
# domain packs under skills/.
install-global: build install-cli
	@profile="$(SPACECRAFT_SKILL_PROFILE)"; \
	if [ "$(FULL)" = "1" ]; then profile=full; fi; \
	if [ -z "$$profile" ]; then profile=lean; fi; \
	export SPACECRAFT_SKILL_PROFILE="$$profile"; \
	sh $(ROOT)/scripts/global-sync.sh "$(ROOT)" "$(GLOBAL)" agents install; \
	echo "agents -> $(GLOBAL)/agents"; \
	sh $(ROOT)/scripts/global-sync.sh "$(ROOT)" "$(GLOBAL)" skills install; \
	echo "skills -> $(GLOBAL)/skills ($$profile)"; \
	test -f "$(GLOBAL)/skills/sc-run/SKILL.md"; \
	test -f "$(GLOBAL)/skills/sc-ship/SKILL.md"; \
	test -f "$(GLOBAL)/skills/sc-quick/SKILL.md"; \
	test -f "$(GLOBAL)/skills/sc-storm/SKILL.md"; \
	test -f "$(GLOBAL)/skills/sc-discuss/references/lens-pass.md"; \
	test -f "$(GLOBAL)/skills/sc-judge/references/judge-break/empty-evidence/expect.json"; \
	python3 $(ROOT)/scripts/mcp-merge.py merge "$(GLOBAL)/mcp.json" "$(ROOT)/.cursor/mcp.json"; \
	sh $(ROOT)/scripts/install-global-hooks.sh "$(ROOT)" "$(GLOBAL)"; \
	sh $(ROOT)/scripts/gen-user-rules.sh "$(ROOT)/.cursor/rules" "$(GLOBAL)/spacecraft/USER-RULES.txt"; \
	echo "user rules -> $(GLOBAL)/spacecraft/USER-RULES.txt (paste into Settings > Rules > User Rules)"; \
	echo "Global install complete. Restart Cursor to pick up /sc-run, /sc-ship, /sc-quick, and sc-storm."

# One-shot machine install: link Node CLI, then clone/update + User-layer via script.
install-machine: install-cli
	@echo "install-cli -> $(LOCAL_BIN)/spacecraft; User-layer via scripts/install-machine.sh"
	@sh $(ROOT)/scripts/install-machine.sh

# Antigravity targets
sync-antigravity:
	node $(ROOT)/scripts/sync-antigravity.mjs

test-antigravity: sync-antigravity
	@sh $(ROOT)/scripts/antigravity-smoke.sh "$(ROOT)"

install-antigravity: install-antigravity-global

install-antigravity-global: build install-cli sync-antigravity
	@sh $(ROOT)/scripts/install-antigravity.sh global
	@echo "Antigravity global plugin installed. Skills and rules active."

install-antigravity-project: build install-cli sync-antigravity
	@sh $(ROOT)/scripts/install-antigravity.sh project "$(PROJECT)"

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
	@if [ -d "$(HOME)/.gemini/config/plugins/spacecraft" ]; then rm -rf "$(HOME)/.gemini/config/plugins/spacecraft" && echo "removed global Antigravity plugin"; fi
	@echo "Uninstall complete. Project-local .cursor/ and .space/ left untouched."

clean:
	rm -f $(BIN)
