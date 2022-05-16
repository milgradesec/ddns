VERSION     := $(shell git describe --tags --always --abbrev=8)
SYSTEM      := 
BUILDFLAGS  := -v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH := github.com/milgradesec/ddns

.PHONY: all
all: build

.PHONY: build
build:
	$(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docker
docker: build
	docker build . -f Dockerfile

.PHONY: release
release:
	docker buildx build . -f build.Dockerfile \
		--platform linux/amd64 \
		--tag ghcr.io/milgradesec/ddns:amd64 \
		--push
	docker buildx build . -f build.Dockerfile \
		--platform linux/arm64 \
		--tag ghcr.io/milgradesec/ddns:arm64 \
		--push
	docker manifest create ghcr.io/milgradesec/ddns:$(VERSION) \
		ghcr.io/milgradesec/ddns:arm64 \
		ghcr.io/milgradesec/ddns:amd64
	docker manifest create ghcr.io/milgradesec/ddns:latest \
		ghcr.io/milgradesec/ddns:arm64 \
		ghcr.io/milgradesec/ddns:amd64
	docker manifest push --purge ghcr.io/milgradesec/ddns:$(VERSION)
	docker manifest push --purge ghcr.io/milgradesec/ddns:latest
