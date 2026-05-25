.PHONY: build webapp dev test vet clean release docker-build docker-up

APP      := featmap
BINDIR   := bin
WEBAPP   := webapp
GOFLAGS  ?= -ldflags="-s -w"
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

## build: build the Go binary (webapp must be built first for go:embed)
build: webapp
	go build $(GOFLAGS) -o $(BINDIR)/$(APP) .

## webapp: build the frontend SPA
webapp:
	cd $(WEBAPP) && npm install && npm run build

## dev: run the server in development mode (uses conf.json)
dev: webapp
	go run .

## test: run all tests
test:
	go test ./...

## vet: run go vet
vet:
	go vet ./...

## clean: remove build artifacts
clean:
	rm -rf $(BINDIR) $(WEBAPP)/build $(WEBAPP)/node_modules

## release: cross-compile for darwin/linux/windows amd64
release: webapp
	@mkdir -p $(BINDIR)
	GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) -o $(BINDIR)/$(APP)-$(VERSION)-darwin-amd64 .
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -o $(BINDIR)/$(APP)-$(VERSION)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BINDIR)/$(APP)-$(VERSION)-windows-amd64.exe .

## docker-build: build the Docker image
docker-build:
	docker build -t $(APP):$(VERSION) .

## docker-up: start services with docker-compose
docker-up:
	FEATMAP_DB=$(FEATMAP_DB) \
	FEATMAP_DB_USER=$(FEATMAP_DB_USER) \
	FEATMAP_DB_PASSWORD=$(FEATMAP_DB_PASSWORD) \
	FEATMAP_HTTP_PORT=$(FEATMAP_HTTP_PORT) \
	docker-compose up -d

## help: show this help
help:
	@grep -E '^## .*: ' $(firstword $(MAKEFILE_LIST)) | sed 's/## //'

# By default, build
.DEFAULT_GOAL := build
