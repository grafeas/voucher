# Builder
FROM golang:1.10-alpine as builder

LABEL maintainer "klautcomputing"
WORKDIR /go/src/github.com/Shopify/voucher
COPY . .
RUN apk update \
&& apk add git \
&& apk add make \
&& make all

# Final build
FROM google/cloud-sdk:alpine

COPY --from=builder /go/src/github.com/Shopify/voucher/voucher_server_unix /usr/local/bin/voucher_server
COPY --from=builder /go/src/github.com/Shopify/voucher/entrypoint.sh /usr/local/entrypoint.sh
COPY --from=builder /go/src/github.com/Shopify/voucher/config/config.toml /etc/voucher/config.toml
COPY config/secrets.production.ejson /etc/voucher/secrets.production.ejson

RUN addgroup -S -g 10000 voucher && \
    adduser -S -u 10000 -G voucher voucher

USER 10000:10000

ENTRYPOINT ["/usr/local/entrypoint.sh"]
