build:
	cd app/ui && yarn build
	go build -o bin/gocan ./app/cmd/gocan/main.go

run:
	go run app/cmd/gocan/main.go

e2e: build
	go test ./e2e/...

internal-test:
	go test ./internal/...

test: build
	go test ./...