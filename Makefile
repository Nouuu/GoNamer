.PHONY: build clean lint format install uninstall test security release-dry-run release-snapshot release release-major release-minor release-patch

BINARY_NAME=gonamer
BUILD_DIR=build
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOLANGCI_LINT_VERSION=v1.62.2

# Version information
VERSION=$(shell git describe --tags 2>/dev/null || git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
CURRENT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# LDFLAGS for build
LDFLAGS=-s -w \
    -X github.com/nouuu/gonamer/cmd.version=$(VERSION) \
	-X github.com/nouuu/gonamer/cmd.commit=$(COMMIT) \
	-X github.com/nouuu/gonamer/cmd.date=$(DATE)

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install github.com/goreleaser/goreleaser/v2@latest

build:
	mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gonamer/main.go
	@echo "Built version $(VERSION) ($(COMMIT))"

install:
	go install -ldflags="$(LDFLAGS)" ./cmd/gonamer
	@echo "Installed version $(VERSION) ($(COMMIT))"

uninstall:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
	rm -f project_knowledge.md
	rm -f *.log
	rm -f *.gob
	rm -f coverage.txt
	rm -rf security-report.json

lint: tools
	golangci-lint run ./...

format:
	gofmt -s -w $(GO_FILES)

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func=coverage.txt

security:
	go list -json -m all | docker run --rm -i sonatypecommunity/nancy:latest sleuth --exclude-vulnerability CVE-2024-45338
	docker run --rm -v $(PWD):/app -w /app securego/gosec:latest -no-fail -fmt=json -out=security-report.json ./...

release-dry-run: tools
	goreleaser release --clean --skip=publish --snapshot

release-snapshot: tools
	goreleaser release --clean --snapshot

# Helper function to create a new release
define do_release
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: Working directory is not clean. Please commit or stash changes first."; \
		exit 1; \
	fi
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "New version: $(1)"
	@read -p "Continue? [y/N] " ans && [ $${ans:-N} = y ]
	@git tag -a $(1) -m "Release $(1)"
	@git push origin $(1)
	@echo "Tag $(1) created and pushed. GitHub Actions will handle the release."
endef

release-major:
	$(eval NEW_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{ printf "v%d.0.0", $$1+1 }'))
	$(call do_release,$(NEW_VERSION))

release-minor:
	$(eval NEW_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{ printf "v%d.%d.0", $$1, $$2+1 }'))
	$(call do_release,$(NEW_VERSION))

release-patch:
	$(eval NEW_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{ printf "v%d.%d.%d", $$1, $$2, $$3+1 }'))
	$(call do_release,$(NEW_VERSION))

help-release:
	@echo "Release targets:"
	@echo "  release-dry-run   - Test release process without publishing"
	@echo "  release-snapshot  - Create a snapshot release for testing"
	@echo "  release-major     - Create a new major release (X.0.0)"
	@echo "  release-minor     - Create a new minor release (0.X.0)"
	@echo "  release-patch     - Create a new patch release (0.0.X)"