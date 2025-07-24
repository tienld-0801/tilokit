.PHONY: all build run run-react run-laravel clean test

all: build

run:
	go run main.go

run-react:
	go run main.go react my-react-app

run-laravel:
	go run main.go laravel my-laravel-app

build:
	go build -o build/tilokit main.go

clean:
	rm -rf build/

test:
	go test ./...
