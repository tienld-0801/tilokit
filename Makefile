.PHONY: all build run run-react run-laravel clean test

all: build

install:
	go mod tidy

run:
	go run . $(ARGS)

run-react:
	go run main.go react my-react-app

run-laravel:
	go run main.go laravel my-laravel-app

build:
	go build -o build/tilokit .

clean:
	rm -rf build/

test:
	go test ./...
