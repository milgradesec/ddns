VERSION:=$(shell git describe --tags --always --abbrev=0)
SYSTEM:=
BUILDFLAGS:=-v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH:=github.com/milgradesec/ddns
CGO_ENABLED:=0

all: build

clean:
	go clean

test:
	go test ./...

lint:
	golangci-lint run

build:
	CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/ddns

docker: 
	docker build . -f build.Dockerfile

release:
	docker buildx build . -f build.Dockerfile \
		--platform linux/amd64 \
		--tag ghcr.io/milgradesec/ddns:amd64 \
		--push
	docker buildx build . -f build.Dockerfile \
		--platform linux/arm64 \
		--tag ghcr.io/milgradesec/ddns:arm64 \
		--push
	docker manifest create ghcr.io/milgradesec/ddns:$(VERSION) ghcr.io/milgradesec/ddns:arm64 ghcr.io/milgradesec/ddns:amd64
	docker manifest create ghcr.io/milgradesec/ddns:latest ghcr.io/milgradesec/ddns:arm64 ghcr.io/milgradesec/ddns:amd64
	docker manifest push --purge ghcr.io/milgradesec/ddns:latest
	docker manifest push --purge ghcr.io/milgradesec/ddns:$(VERSION)
