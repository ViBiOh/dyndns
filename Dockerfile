FROM scratch

ENTRYPOINT [ "/dyndns" ]

ARG VERSION
ENV VERSION=${VERSION}

ARG TARGETOS
ARG TARGETARCH

COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY release/dyndns_${TARGETOS}_${TARGETARCH} /dyndns
