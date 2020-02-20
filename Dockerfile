FROM debian:stable-slim

ENTRYPOINT [ "/dyndns" ]

ARG VERSION
ENV VERSION=${VERSION}

COPY dyndns /dyndns
