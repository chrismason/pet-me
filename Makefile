TOOLS=$(CURDIR)/_tools
export GO111MODULE=on
export CGO_ENABLED=0

.PHONY: all
all: build

.PHONY: build
build: api cli

.PHONY: tools
tools:
	mkdir -p $(TOOLS)
	@[ ! -f $(TOOLS)/go.mod ] && cd $(TOOLS) && go mod init tools || true
	cd $(TOOLS) && go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.24.0

.PHONY: clean
clean:
	go clean -cache -testcache -modcache

.PHONY: api
api:
	go build -o bin/pet-me-api ./cmd/api

.PHONY: cli
cli:
	go build -o bin/pet-me ./cmd/cli

.PHONY: fmt
fmt:
	@echo "==> running Go Format"
	gofmt -s -l -w .

.PHONY: lint
lint:
	@echo "==> linting Go code"
	@$(shell go env GOPATH)/bin/golangci-lint run ./...
	@echo "==> running go vet"
	go vet ./...

.PHONY: test
test:
	@echo "==> running Go tests"
	CGO_ENABLED=1 go test -p 1 -race ./...
