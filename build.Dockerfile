FROM --platform=amd64 golang:1.16.6

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/app
COPY . .

RUN make build SYSTEM="GOOS=${TARGETOS} GOARCH=${TARGETARCH}"

FROM alpine:3.14.0

RUN apk update && apk add --no-cache ca-certificates 

FROM scratch

COPY --from=0 /go/src/app/ddns /ddns
COPY --from=1 /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/ddns"]
