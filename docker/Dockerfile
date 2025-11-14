FROM quay.io/prometheus/busybox-linux-amd64:latest
LABEL maintainer="Vectinx <vectinx@yandex.ru>"

COPY ./389-ds-exporter /bin/389-ds-exporter
RUN chmod 700 /bin/389-ds-exporter

EXPOSE      9389
USER        nobody
ENTRYPOINT  [ "/bin/389-ds-exporter" ]
