.PHONY: build test clean lint install uninstall dev docs

GO_SRC := .engine/scripts/src
GO_OUT := scripts/spacecraft
LOCAL_BIN := $(HOME)/.local/bin
GLOBAL_CONFIG := $(HOME)/.config/opencode/opencode.jsonc
SPACECRAFT_ROOT := $(CURDIR)

build:
	cd $(GO_SRC) && go build -o ../../../$(GO_OUT) .

test:
	cd $(GO_SRC) && go test -coverprofile=coverage.out ./...

clean:
	rm -f $(GO_OUT)

lint:
	cd $(GO_SRC) && go vet ./...

install: build
	@mkdir -p $(LOCAL_BIN)
	@ln -sf $(SPACECRAFT_ROOT)/$(GO_OUT) $(LOCAL_BIN)/spacecraft
	@mkdir -p $(HOME)/.config/opencode
	@printf '{\n  "$$schema": "https://opencode.ai/config.json",\n  "default_agent": "sc-commander",\n  "share": "disabled",\n  "lsp": true,\n  "instructions": ["%s/.engine/PERSONA.md", "%s/.engine/AGENTS.md"],\n  "skills": {"paths": ["%s/.engine/skills"]},\n  "plugin": ["%s/.opencode/plugins/engine.js"],\n  "agent": {\n    "sc-commander": {\n      "mode": "primary",\n      "description": "Mission commander — orchestrates mission-driven implementation using lean prompts",\n      "model": "opencode-go/deepseek-v4-pro",\n      "temperature": 0.2,\n      "variant": "high",\n      "color": "#4fc3f7"\n    },\n    "sc-planner": {\n      "mode": "subagent",\n      "description": "Read-only planner that turns a mission spec into a small executable plan",\n      "model": "opencode-go/kimi-k2.7-code",\n      "temperature": 0.1,\n      "variant": "low",\n      "color": "#80cbc4"\n    },\n    "sc-designer": {\n      "mode": "subagent",\n      "description": "Read-only design agent for UI direction, critique, and anti-slop review",\n      "model": "opencode-go/glm-5.2",\n      "temperature": 0.3,\n      "variant": "medium",\n      "color": "#ce93d8"\n    },\n    "sc-reviewer": {\n      "mode": "subagent",\n      "description": "Read-only reviewer for diff, evidence, and release readiness",\n      "model": "opencode-go/deepseek-v4-pro",\n      "temperature": 0.1,\n      "variant": "medium",\n      "color": "#ef9a9a"\n    },\n    "sc-coder": {\n      "mode": "subagent",\n      "description": "Write-capable coder that implements production code to satisfy tasks and tests",\n      "model": "opencode-go/kimi-k2.7-code",\n      "temperature": 0.2,\n      "variant": "low",\n      "color": "#a5d6a7"\n    },\n    "sc-tester": {\n      "mode": "subagent",\n      "description": "Write-capable tester that writes tests and captures verification evidence (TDD)",\n      "model": "opencode-go/deepseek-v4-flash",\n      "temperature": 0.1,\n      "color": "#ffcc80"\n    }\n  }\n}\n' $(SPACECRAFT_ROOT) $(SPACECRAFT_ROOT) $(SPACECRAFT_ROOT) $(SPACECRAFT_ROOT) > $(GLOBAL_CONFIG)
	@echo "Installed: spacecraft -> $(LOCAL_BIN)/spacecraft"
	@echo "Config:    $(GLOBAL_CONFIG)"
	@echo "Restart OpenCode to apply."

uninstall:
	@rm -f $(LOCAL_BIN)/spacecraft
	@echo "Removed $(LOCAL_BIN)/spacecraft"
	@echo "Config at $(GLOBAL_CONFIG) left intact."

dev: build
	@echo "Development mode: binary rebuilt at $(GO_OUT)"
	@echo "Run: scripts/spacecraft help"

docs:
	@echo "Generating documentation..."
	@mkdir -p docs/generated
	@scripts/spacecraft help > docs/generated/cli-reference.txt 2>&1
	@echo "Documentation generated in docs/generated/"
