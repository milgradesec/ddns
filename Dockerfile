FROM gcr.io/distroless/static-debian11:nonroot

ADD ddns /ddns

USER nonroot
ENTRYPOINT ["/ddns"]
