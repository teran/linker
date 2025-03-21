FROM alpine:3.20.3 AS certificates

RUN apk add --update --no-cache \
  ca-certificates=20241121-r1

FROM scratch

COPY dockerfiles/rootfs/etc/passwd /etc/passwd
COPY dockerfiles/rootfs/etc/group /etc/group

COPY --from=certificates /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --chmod=0755 --chown=root:root dist/migrator_linux_amd64_v3/migrator /migrator
COPY --chmod=0644 --chown=root:root repositories/metadata/postgresql/migrations/sql /migrations/sql

USER nobody

ENTRYPOINT [ "/migrator" ]
