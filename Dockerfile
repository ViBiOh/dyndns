FROM scratch

ENV ZONEINFO zoneinfo.zip
EXPOSE 1080

HEALTHCHECK --retries=5 CMD [ "/goweb", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/goweb" ]

ARG VERSION
ENV VERSION=${VERSION}

ARG TARGETOS
ARG TARGETARCH

COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY zoneinfo.zip /
COPY release/goweb_${TARGETOS}_${TARGETARCH} /goweb
