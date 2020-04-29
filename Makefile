VERSION:=$(shell git describe --tags --always --dirty="-dev")
BUILDFLAGS:=-v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH:=github.com/milgradesec/ddns
DOCKER_PLATFORM:=linux/arm/v7

.PHONY: all
all: build

.PHONY: build
build:
	go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns

.PHONY: docker
docker:
	docker build . -t ddns:$(VERSION)

.PHONY: release
release:
	docker buildx build --platform $(DOCKER_PLATFORM) --build-arg VERSION=$(VERSION) . -t milgradesec/ddns:latest -t milgradesec/ddns:$(VERSION) --push

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean