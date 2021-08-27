FROM --platform=amd64 golang:1.17.0

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/app
COPY . .

RUN make build SYSTEM="GOOS=${TARGETOS} GOARCH=${TARGETARCH}"

FROM alpine:3.14.2

RUN apk --update --no-cache add ca-certificates \
    && addgroup -g 1000 ddns \
    && adduser -u 1000 -G ddns -s /sbin/nologin -D ddns

FROM scratch

COPY --from=0 /go/src/app/ddns /ddns
COPY --from=1 /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/ddns"]
