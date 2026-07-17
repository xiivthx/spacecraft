.PHONY: build install uninstall clean

CURSOR_AGENTS := $(HOME)/.cursor/agents
CURSOR_RULES := $(HOME)/.cursorrules
CURSOR_MCP   := $(HOME)/.cursor/mcp.json
LOCAL_BIN    := $(HOME)/.local/bin
ROOT         := $(CURDIR)

build:
	cd $(ROOT)/cmd/spacecraft && go build -o $(ROOT)/spacecraft .

install: build
	@mkdir -p $(CURSOR_AGENTS) $(LOCAL_BIN)
	@cp $(ROOT)/.cursor/agents/sc-*.md $(CURSOR_AGENTS)/
	@cat $(ROOT)/.cursor/rules/050-style.mdc $(ROOT)/.cursor/rules/000-spacecraft.mdc $(ROOT)/.cursor/rules/100-conventions.mdc > $(CURSOR_RULES)
	@cp $(ROOT)/.cursor/mcp.json $(CURSOR_MCP)
	@ln -sf $(ROOT)/spacecraft $(LOCAL_BIN)/spacecraft
	@echo "Installed:"
	@echo "  agents  -> $(CURSOR_AGENTS)"
	@echo "  rules   -> $(CURSOR_RULES)"
	@echo "  mcp     -> $(CURSOR_MCP)"
	@echo "  cli     -> $(LOCAL_BIN)/spacecraft"

uninstall:
	@rm -f $(CURSOR_AGENTS)/sc-*.md
	@rm -f $(CURSOR_RULES)
	@rm -f $(CURSOR_MCP)
	@rm -f $(LOCAL_BIN)/spacecraft
	@echo "Removed spacecraft from $(HOME)/.cursor and $(LOCAL_BIN)"

clean:
	rm -f $(ROOT)/spacecraft
