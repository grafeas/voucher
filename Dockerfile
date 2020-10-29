# Builder
FROM golang:1.14-alpine as builder

LABEL maintainer "catherinejones"
WORKDIR /go/src/github.com/grafeas/voucher
COPY . .
RUN apk --no-cache add \
    git \
    make && \
    make voucher_server

# Final build
FROM alpine:3.12

COPY --from=builder /go/src/github.com/grafeas/voucher/build/voucher_server /usr/local/bin/voucher_server
COPY --from=builder /go/src/github.com/grafeas/voucher/entrypoint.sh /usr/local/entrypoint.sh
COPY --from=builder /go/src/github.com/grafeas/voucher/config/config.toml /etc/voucher/config.toml
COPY config/secrets.production.ejson /etc/voucher/secrets.production.ejson

RUN apk add --no-cache \
    ca-certificates && \
    addgroup -S -g 10000 voucher && \
    adduser -S -u 10000 -G voucher voucher

USER 10000:10000

ENTRYPOINT ["/usr/local/entrypoint.sh"]
