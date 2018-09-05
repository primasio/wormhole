# Makefile for wormhole distribution package

all: dist

.PHONY: deps dist
	dist: build-linux-x64
	mkdir dist/config
	cp config/*.yaml dist/config/
	cp index.html dist/

.PHONY: test
test: deps
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: clean
clean:
	go clean
	rm -rf dist/*

.PHONY: deps
deps:
	go get ./...

.PHONY: build-linux-x64
build-linux-x64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/wormhole -v
