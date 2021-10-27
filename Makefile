SHELL := /bin/bash
VERSION := 0.2.2
PLATFORMS := linux/amd64 darwin/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

frontend:
	cd app/ui && yarn install && CI=true yarn build

backend:
	go build -ldflags="-X 'main.Version=v$(VERSION)'" -o bin/gocan ./app/cmd/gocan/main.go

doc: backend
	rm -rf ./doc/commands
	mkdir -p ./doc/commands
	./bin/gocan generate-doc --directory ./doc/commands

build: frontend backend doc

run:
	go run app/cmd/gocan/main.go

e2e: build
	go test ./e2e/...

business-test:
	go test ./business/...

test: build
	go test ./...

release: frontend $(PLATFORMS)


docker: frontend backend
	docker build \
		-f docker/Dockerfile  \
		-t gocan:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -ldflags="-X 'main.Version=v$(VERSION)'" -o 'bin/gocan-$(os)-$(arch)' ./app/cmd/gocan/main.go

.PHONY: release $(PLATFORMS)