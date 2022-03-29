FROM alpine:3.15.3

RUN apk update && \
    apk upgrade --available && \
    apk add --no-cache ca-certificates && \
    addgroup -S ddns && \
    adduser -S -G ddns ddns

FROM scratch

COPY --from=0 /etc/ssl/certs /etc/ssl/certs
COPY --from=0 /etc/passwd /etc/passwd

ADD ddns /ddns

USER ddns
ENTRYPOINT ["/ddns"]
