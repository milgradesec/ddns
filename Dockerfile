FROM --platform=linux/amd64 golang:1.14.2-alpine AS builder

WORKDIR /go/src/app
COPY . .

ARG VERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -v -ldflags "-s -w -X main.Version=${VERSION}" github.com/milgradesec/ddns/cmd/ddns

FROM alpine:3.11.3

RUN apk update && apk add --no-cache ca-certificates && \
    addgroup -S ddns && adduser -S -G ddns ddns

COPY --from=0 /go/src/app/ddns /ddns
USER ddns
ENTRYPOINT ["/ddns"]