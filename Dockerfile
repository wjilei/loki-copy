FROM alpine:latest

COPY loki-copy /app/loki-copy
COPY config.yaml /app/config.yaml

WORKDIR /app

CMD ["./loki-copy"]