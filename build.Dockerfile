FROM golang:1.14.6

WORKDIR /go/src/app
COPY . .

RUN make build

FROM alpine:3.12

RUN apk update && apk add --no-cache ca-certificates 

FROM scratch

COPY --from=0 /go/src/app/ddns /ddns
COPY --from=1 /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/ddns"]