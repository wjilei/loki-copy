FROM alpine:latest

COPY loki-copy /app/loki-copy
COPY config.yaml /app/config.yaml

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /app

CMD ["./loki-copy"]