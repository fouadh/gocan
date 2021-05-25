build:
	go build -o bin/gocan ./cmd/gocan/main.go

run:
	go run cmd/gocan/main.go

it:
	go test cmd/gocan/tests/*.go