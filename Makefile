VERSION:=$(shell git describe --tags --always --abbrev=0 --dirty="-dev")
SYSTEM:=
BUILDFLAGS:=-v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH:=github.com/milgradesec/ddns
CGO_ENABLED:=0

.PHONY: all
all: build

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns

.PHONY: docker
.ONESHELL:
docker:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns
	docker.exe build . -t ddns:$(VERSION)

.PHONY: release
.ONESHELL:
release:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm64 go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns
	docker.exe buildx build --platform linux/arm64 . -t milgradesec/ddns:latest -t milgradesec/ddns:$(VERSION) --push

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run