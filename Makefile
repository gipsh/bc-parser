SHELL := /bin/bash

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

build:
	go build -o parser-api cmd/api/main.go

	
test:
	go test -v -cover ./...

run:
	go run main.go


.PHONY: build test run