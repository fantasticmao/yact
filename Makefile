NAME=yact
BIN_DIR=bin
VERSION=$(shell cat VERSION)

.PHONY: get format test run build

get:
	go get -u -t ./...

format:
	go fmt ./...

test:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic -v

run:
	go run main.go

build:
	go build -o $(BIN_DIR)/$(NAME)
