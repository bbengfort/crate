# Shell to use with Make
SHELL := /bin/sh

deps:
	@echo "Fetching Dependencies"
	@echo "Not implemented yet"

fmt:
	@echo "Formatting"
	gofmt -w .

test:
	ginkgo -r

.PHONY:
	deps fmt test
