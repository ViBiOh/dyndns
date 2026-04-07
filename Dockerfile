FROM rg.fr-par.scw.cloud/vibioh/scratch

ENTRYPOINT [ "/dyndns" ]

ARG VERSION
ENV VERSION=${VERSION}

ARG GIT_SHA
ENV GIT_SHA=${GIT_SHA}

ARG TARGETOS
ARG TARGETARCH

COPY cacert.pem /etc/ssl/cert.pem
COPY release/dyndns_${TARGETOS}_${TARGETARCH} /dyndns
