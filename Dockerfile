FROM alpine:3.13

RUN apk update && apk add --no-cache ca-certificates 

FROM scratch

COPY --from=0 /etc/ssl/certs /etc/ssl/certs

ADD ddns /ddns
ENTRYPOINT ["/ddns"]
