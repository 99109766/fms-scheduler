## run: Run the application with the config file
.PHONY: run
run:
	@echo "\033[0;32mRunning the application...\033[0m"
	go run ./cmd/main.go --config config.yaml

## build: Build the application binary
.PHONY: build
build:
	@echo "\033[0;36mBuilding the application...\033[0m"
	go build -o ./bin/main ./cmd/main.go

## test: Run tests to check code validity
.PHONY: test
test:
	@echo "\033[0;33mRunning tests...\033[0m"
	go test -v ./...

## lint: Run linter to check code quality
.PHONY: lint
lint:
	@echo "\033[0;35mRunning linter...\033[0m"
	golangci-lint run

## install: Install dependencies
.PHONY: install
install:
	@echo "\033[0;34mInstalling dependencies...\033[0m"
	go mod tidy

## clean: Clean up the project binaries
.PHONY: clean
clean:
	@echo "\033[0;30mCleaning up...\033[0m"
	rm -rf ./bin

## help: Show help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.DEFAULT_GOAL := run