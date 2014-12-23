# Shell to use with Make
SHELL := /bin/sh

all: fmt deps
	@echo "Building"
	@mkdir -p bin/
	@go build -v -o bin/crate .

deps:
	@echo "Fetching Dependencies"
	@go get -d -v ./crate/...

fmt:
	@echo "Formatting"
	gofmt -w .

test: deps
	ginkgo -r -v

.PHONY:
	all deps fmt test
