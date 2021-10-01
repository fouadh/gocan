build:
	cd app/ui && yarn install && yarn build
	go build -o bin/gocan ./app/cmd/gocan/main.go

run:
	go run app/cmd/gocan/main.go

e2e: build
	go test ./e2e/...

business-test:
	go test ./business/...

test: build
	go test ./...