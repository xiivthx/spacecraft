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
