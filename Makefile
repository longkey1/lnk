.DEFAULT_GOAL := help

export GO_VERSION=$(shell if [ -f go.mod ]; then grep "^go " go.mod | sed 's/^go //'; else curl -s https://go.dev/dl/ | grep -o 'go[0-9]\+\.[0-9]\+' | head -1 | sed 's/go//'; fi)
export PRODUCT_NAME := lnkr

.PHONY: init
init: ## Initialize the project
	mkdir -p .devcontainer
	cd .devcontainer && cat devcontainer.json.dist | envsubst '$${GO_VERSION} $${PRODUCT_NAME}' > devcontainer.json

.PHONY: build
build: ## Build the application
	go build -o bin/lnkr .

.PHONY: build-dev
build-dev: ## Build the application with development version info
	@COMMIT_SHA=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	BUILD_TIME=$$(date -u '+%Y-%m-%dT%H:%M:%SZ'); \
	LDFLAGS="-X github.com/longkey1/lnkr/internal/version.Version=dev -X github.com/longkey1/lnkr/internal/version.CommitSHA=$$COMMIT_SHA -X github.com/longkey1/lnkr/internal/version.BuildTime=$$BUILD_TIME"; \
	go build -ldflags "$$LDFLAGS" -o bin/lnkr .

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf dist/

# Get current version from git tag
CURRENT_VERSION := $(shell git tag --sort=-v:refname | head -n1 2>/dev/null | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' || echo "v0.0.0")
VERSION := $(CURRENT_VERSION)
MAJOR := $(shell echo $(VERSION) | cut -d. -f1 | tr -d 'v' | sed 's/^$$/0/')
MINOR := $(shell echo $(VERSION) | cut -d. -f2 | sed 's/^$$/0/')
PATCH := $(shell echo $(VERSION) | cut -d. -f3 | sed 's/^$$/0/')

# Variables for release target
dryrun ?= true
type ?=

.PHONY: release
release: ## Release target with type argument. Usage: make release type=patch|minor|major dryrun=false
	@if [ "$(type)" = "" ]; then \
		echo "Usage: make release type=<type> [dryrun=false]"; \
		echo ""; \
		echo "Types:"; \
		echo "  patch  - Increment patch version (e.g., v1.2.3 -> v1.2.4)"; \
		echo "  minor  - Increment minor version (e.g., v1.2.3 -> v1.3.0)"; \
		echo "  major  - Increment major version (e.g., v1.2.3 -> v2.0.0)"; \
		echo ""; \
		echo "Options:"; \
		echo "  dryrun - Set to false to actually create and push the tag (default: true)"; \
		echo ""; \
		echo "Current version: $(CURRENT_VERSION)"; \
		exit 0; \
	elif [ "$(type)" = "patch" ] || [ "$(type)" = "minor" ] || [ "$(type)" = "major" ]; then \
		NEXT_VERSION=$$(if [ "$(type)" = "patch" ]; then \
			echo "v$(MAJOR).$(MINOR).$$(($(PATCH) + 1))"; \
		elif [ "$(type)" = "minor" ]; then \
			echo "v$(MAJOR).$$(($(MINOR) + 1)).0"; \
		elif [ "$(type)" = "major" ]; then \
			echo "v$$(($(MAJOR) + 1)).0.0"; \
		fi); \
		echo "Current version: $(CURRENT_VERSION)"; \
		echo "Next version: $$NEXT_VERSION"; \
		if [ "$(dryrun)" = "false" ]; then \
			echo "Creating new tag $$NEXT_VERSION..."; \
			git push origin master --no-verify --force-with-lease; \
			git tag -a $$NEXT_VERSION -m "Release of $$NEXT_VERSION"; \
			git push origin $$NEXT_VERSION --no-verify --force-with-lease; \
			echo "Tag $$NEXT_VERSION has been created and pushed"; \
			echo "Running goreleaser to build and release..."; \
			goreleaser release --clean; \
		else \
			echo "[DRY RUN] Showing what would be done..."; \
			echo "Would push to origin/master"; \
			echo "Would create tag: $$NEXT_VERSION"; \
			echo "Would push tag to origin: $$NEXT_VERSION"; \
			echo "Would run goreleaser release"; \
			echo ""; \
			echo "To execute this release, run:"; \
			echo "  make release type=$(type) dryrun=false"; \
			echo "Dry run complete."; \
		fi \
	else \
		echo "Error: Invalid release type. Use 'patch', 'minor', or 'major'"; \
		exit 1; \
	fi

.PHONY: re-release

# Variables for re-release target
dryrun ?= true
tag ?=

re-release: ## Rerelease target with tag argument. Usage: make re-release tag=<tag> dryrun=false
	@TAG="$(tag)"; \
	if [ -z "$$TAG" ]; then \
		TAG=$$(git describe --tags --abbrev=0); \
	fi; \
	if [ -z "$$TAG" ]; then \
		echo "Error: No tag found near HEAD and no tag specified."; \
		exit 1; \
	fi; \
	echo "Target tag: $$TAG"; \
	if [ "$(dryrun)" = "false" ]; then \
		echo "Deleting GitHub release..."; \
		gh release delete "$$TAG" -y; \
		echo "Deleting local tag..."; \
		git tag -d "$$TAG"; \
		echo "Deleting remote tag..."; \
		git push origin ":refs/tags/$$TAG" --no-verify --force; \
		echo "Recreating tag on HEAD..."; \
		git tag -a "$$TAG" -m "Release $$TAG"; \
		echo "Pushing tag to origin..."; \
		git push origin "$$TAG" --no-verify --force-with-lease; \
		echo "Recreating GitHub release with goreleaser..."; \
		goreleaser release --clean; \
		echo "Done!"; \
	else \
		echo "[DRY RUN] Showing what would be done..."; \
		echo "Would delete release: $$TAG"; \
		echo "Would delete local tag: $$TAG"; \
		echo "Would delete remote tag: $$TAG"; \
		echo "Would create new tag at HEAD: $$TAG"; \
		echo "Would push tag to origin: $$TAG"; \
		echo "Would run goreleaser release"; \
		echo ""; \
		echo "To execute this re-release, run:"; \
		if [ -n "$(tag)" ]; then \
			echo "  make re-release tag=$$TAG dryrun=false"; \
		else \
			echo "  make re-release dryrun=false"; \
		fi; \
		echo "Dry run complete."; \
	fi

.PHONY: release-dry-run
release-dry-run: ## Run goreleaser in dry-run mode
	goreleaser release --snapshot --clean --skip-publish

.PHONY: release-snapshot
release-snapshot: ## Create a snapshot release
	goreleaser release --snapshot --clean --skip-publish

.PHONY: install
install: ## Install goreleaser
	go install github.com/goreleaser/goreleaser@latest

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


