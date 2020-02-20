FROM alpine

ENTRYPOINT [ "/dyndns" ]

ARG VERSION
ENV VERSION=${VERSION}

COPY dyndns /dyndns
