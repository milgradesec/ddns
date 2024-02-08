FROM gcr.io/distroless/static-debian12:nonroot

ADD ddns /ddns

USER nonroot
ENTRYPOINT ["/ddns"]
