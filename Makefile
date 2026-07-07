.PHONY: help test build spacecraft \
        sc-init sc-status sc-missions sc-use sc-flow \
        sc-git sc-git-suggest sc-archive sc-closeout sc-validate \
        sc-design-serve sc-design-open

SCRIPTS_BIN := scripts/spacecraft

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: | build ## Run tests (build required first for Go binary)
	node --test tests/*.test.mjs

build: ## Build the Go binary
	cd scripts/src && go build -o ../spacecraft .

spacecraft: ## Run the spacecraft CLI
	$(SCRIPTS_BIN)

sc-init: ## Initialize spacecraft
	$(SCRIPTS_BIN) init

sc-status: ## Print spacecraft status
	$(SCRIPTS_BIN) status

sc-missions: ## List missions
	$(SCRIPTS_BIN) missions

sc-use: ## Select a mission (use ARGS="<selector>")
	$(SCRIPTS_BIN) use $(ARGS)

sc-flow: ## Check workflow readiness
	$(SCRIPTS_BIN) flow

sc-git: ## Print git safety status
	$(SCRIPTS_BIN) git-info

sc-git-suggest: ## Suggest branch and commit names
	$(SCRIPTS_BIN) git-suggest

sc-archive: ## Archive a shipped mission
	$(SCRIPTS_BIN) archive

sc-closeout: ## Check release closeout readiness
	$(SCRIPTS_BIN) closeout-check

sc-validate: ## Validate resolved mission artifacts
	$(SCRIPTS_BIN) validate

sc-design-serve: ## Serve design HTML artifacts (use ARGS="[artifact-or-dir]")
	node .opencode/skills/sc-design/scripts/serve-html.mjs $(ARGS)

sc-design-open: ## Serve and open design HTML artifacts in browser (use ARGS="[artifact-or-dir]")
	node .opencode/skills/sc-design/scripts/serve-html.mjs --open $(ARGS)
