# Go Task Tracker - Makefile

APP_NAME = task-cli
GO_FILES = ./...

.PHONY: all build clean test lint run

all: clean build

build:
	go build -o $(APP_NAME) main.go

clean:
	rm -f $(APP_NAME)

test:
	go test $(GO_FILES) -v

lint:
	golangci-lint run || true

run:
	go run main.go