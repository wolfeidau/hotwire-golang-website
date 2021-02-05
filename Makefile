MODULE_PKG := github.com/wolfeidau/hotwire-golang-website
WATCH := (.go$$)|(.html$$)

GOLANGCI_VERSION = 1.31.0

GIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y%m%dT%H%M%S')

# This path is used to cache binaries used for development and can be overridden to avoid issues with osx vs linux
# binaries.
BIN_DIR ?= $(shell pwd)/bin

default: clean build archive deploy-bucket package deploy-api
.PHONY: default

deploy: build archive package deploy-api
.PHONY: deploy

ci: clean lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w -X $(MODULE_PKG)/internal/app.BuildDate=${BUILD_DATE} -X $(MODULE_PKG)/internal/app.Commit=${GIT_HASH}"

$(BIN_DIR)/golangci-lint: $(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} $(BIN_DIR)/golangci-lint
$(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv $(BIN_DIR)/golangci-lint $@

$(BIN_DIR)/mockgen:
	@go get -u github.com/golang/mock/mockgen
	@env GOBIN=$(BIN_DIR) GO111MODULE=on go install github.com/golang/mock/mockgen

$(BIN_DIR)/gosec:
	@go get -u github.com/securego/gosec/v2/cmd/gosec@master
	@env GOBIN=$(BIN_DIR) GO111MODULE=on go install github.com/securego/gosec/v2/cmd/gosec

$(BIN_DIR)/reflex:
	@go get -u github.com/cespare/reflex
	@env GOBIN=$(BIN_DIR) GO111MODULE=on go install github.com/cespare/reflex

mocks: $(BIN_DIR)/mockgen
	@echo "--- build all the mocks"
	@bin/mockgen -destination=mocks/session_store.go -package=mocks github.com/dghubble/sessions Store
.PHONY: mocks

clean:
	@echo "--- clean all the things"
	@rm -rf ./dist
.PHONY: clean

scanpr: $(BIN_DIR)/gosec
	$(BIN_DIR)/gosec -fmt golint ./...

scan-report: $(BIN_DIR)/gosec
	$(BIN_DIR)/gosec -no-fail -fmt sarif -out results.sarif ./...
.PHONY: scan-report

lint: $(BIN_DIR)/golangci-lint
	@echo "--- lint all the things"
	@$(BIN_DIR)/golangci-lint run
.PHONY: lint

lint-fix: $(BIN_DIR)/golangci-lint
	@echo "--- lint all the things"
	@$(BIN_DIR)/golangci-lint run --fix
.PHONY: lint-fix

test:
	@echo "--- test all the things"
	@go test -coverprofile=coverage.txt ./...
	@go tool cover -func=coverage.txt
.PHONY: test

install:
	@cd assets && npm ci
.PHONY: install

watch: $(BIN_DIR)/reflex install
	$(BIN_DIR)/reflex -r "$(WATCH)" -s -- make start
.PHONY: watch

start:
	go build $(LDFLAGS) -trimpath -o dist/hotwire-server ./cmd/hotwire-server/main.go
	LOCAL=true ./dist/hotwire-server
.PHONY: start

build:
	@echo "--- build all the things"
	@mkdir -p dist
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -trimpath -o dist ./cmd/...
.PHONY: build

certs:
	@mkdir -p .certs
	@mkcert -cert-file .certs/hotwire.localhost.pem -key-file .certs/hotwire.localhost.key hotwire.localhost
.PHONY: certs
