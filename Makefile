PLATFORMS := linux/amd64 darwin/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

frontend:
	cd app/ui && yarn install && yarn build

build: frontend

	go build -o bin/gocan ./app/cmd/gocan/main.go

run:
	go run app/cmd/gocan/main.go

e2e: build
	go test ./e2e/...

business-test:
	go test ./business/...

test: build
	go test ./...

release: frontend $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o 'bin/gocan-$(os)-$(arch)' ./app/cmd/gocan/main.go

.PHONY: release $(PLATFORMS)