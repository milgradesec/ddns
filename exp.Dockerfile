FROM golang:1.14.4-alpine

WORKDIR /go/src/app
RUN make build

FROM alpine:3.12

RUN apk update && apk add --no-cache ca-certificates 

FROM scratch

COPY --from=1 /etc/ssl/certs /etc/ssl/certs
COPY --from=0 /go/src/app/ddns /ddns

ENTRYPOINT ["/ddns"]