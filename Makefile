VERSION:=$(shell git describe --tags --always --abbrev=0 --dirty="-dev")
BUILDFLAGS:=-v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH:=github.com/milgradesec/ddns
DOCKER_PLATFORM:=linux/arm/v7
SYSTEM:=

ifeq ($(SYSTEM),)
endif

.PHONY: all
all: build

.PHONY: build
build:
	$(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns

.PHONY: docker
.ONESHELL:
docker:
	set CGO_ENABLED=0
	set GOOS=linux
	set GOARCH=amd64
	go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns
	docker build . -t ddns:$(VERSION)

.PHONY: release
.ONESHELL:
release:
	set CGO_ENABLED=0
	set GOOS=linux
	set GOARCH=arm
	set GOARM=7
	go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns
	docker buildx build --platform $(DOCKER_PLATFORM) . -t milgradesec/ddns:latest -t milgradesec/ddns:$(VERSION) --push

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	go clean