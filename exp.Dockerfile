FROM golang:1.14.4

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 go build -v -ldflags="-s -w -X main.Version=DEV" github.com/milgradesec/ddns/cmd/ddns

FROM alpine:3.12

RUN apk update && apk add --no-cache ca-certificates 

FROM scratch

COPY --from=0 /go/src/app/ddns /ddns
COPY --from=1 /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/ddns"]