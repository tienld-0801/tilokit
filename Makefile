run:
	go run main.go

run-react:
	go run main.go react my-react-app

run-laravel:
	go run main.go laravel my-laravel-app

build:
	go build -o build/tilokit main.go
