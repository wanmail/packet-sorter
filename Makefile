# Makefile for compiling Go project for x64 and ARM architectures

BINARY_NAME=packet-sorter
SRC_DIR=./cmd/main.go

.PHONY: all clean build

all: clean build

clean:
	rm -f $(BINARY_NAME)-*

build: build-x64 build-arm

build-x64:
	GOARCH=amd64 GOOS=linux go build -o $(BINARY_NAME)-x64 $(SRC_DIR)

build-arm:
	GOARCH=arm64 go build -o $(BINARY_NAME)-arm $(SRC_DIR)