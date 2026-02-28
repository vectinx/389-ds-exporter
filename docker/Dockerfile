FROM quay.io/prometheus/busybox-linux-amd64@sha256:86d2c6d3d8550a58ab6b61eda6343e523785bf9e7a99fd5efd3d3159e4a2cfb9

LABEL org.opencontainers.image.authors="Vectinx <vectinx@yandex.ru>"
LABEL org.opencontainers.image.source="https://github.com/vectinx/389-ds-exporter"
LABEL org.opencontainers.image.description="Prometheus exporter for 389 Directory Server"
LABEL org.opencontainers.image.licenses="MIT"

COPY --chown=65534:65534 ./389-ds-exporter /bin/389-ds-exporter

EXPOSE      9389
USER 65534:65534
ENTRYPOINT  [ "/bin/389-ds-exporter" ]
