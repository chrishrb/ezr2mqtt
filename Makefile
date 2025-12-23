# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

GOCMD=go
LDFLAGS="-s -w ${LDFLAGS_OPT}"

all: vendor build format lint e2e ## Format, lint, build, test

.PHONY: run
run: ## Run
	go run -tags debug main.go $(RUN_ARGS)

.PHONY: build
build: ## Build
	mkdir -p bin
	go build -tags debug -o bin/ezr2mqtt main.go

.PHONY: vendor
vendor: ## Vendor
	go mod vendor

.PHONY: test
test: ## Test
	${GOCMD} test ./...

.PHONY: e2e
e2e: ## Run unit and e2e tests
	${GOCMD} test -tags=e2e ./...

.PHONY: compile
compile: ## Compile for every OS and Platform
	echo "Compiling for every OS and Platform"
	GOOS=darwin GOARCH=amd64 go build -o bin/ezr2mqtt-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/ezr2mqtt-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/ezr2mqtt-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/ezr2mqtt-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/ezr2mqtt-windows-amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build -o bin/ezr2mqtt-windows-arm64.exe main.go

.PHONY: format
format: ## Format code
	${GOCMD} fmt ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: clean
clean: ## Cleanup build dir
	rm -r bin/

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
