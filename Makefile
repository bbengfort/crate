# Shell to use with Make
SHELL := /bin/sh

deps:
	@echo "Fetching Dependencies"
	@echo "Not implemented yet"

fmt:
	@echo "Formatting"
	gofmt -w .

test:
	go test -v ./crate/...
	- go test -v ./version/...

.PHONY:
	deps fmt test
