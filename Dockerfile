# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: MPL-2.0



FROM golang:1.20.2-alpine3.17 as builder
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/app && \
adduser -u 1002 -D appuser appuser

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/spectrocloud/hello-universe-api"
LABEL org.opencontainers.image.description "A Spectro Cloud demo application intended for learning and showcasing products. This is the API server for Hello Universe."
ENV PORT 3000
ENV HOST '0.0.0.0'
ENV DB_HOST '0.0.0.0'
ENV DB_PORT 5432
ENV DB_USER postgres
ENV DB_PASSWORD password
ENV DB_NAME counter
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=appuser:appuser /go/bin/app /usr/bin/app
RUN apk -U upgrade --no-cache && apk add --no-cache tzdata bash curl openssl jq bind-tools && \
chmod a+x /usr/bin/app
USER appuser
EXPOSE 3000
EXPOSE 5432
CMD ["/usr/bin/app"]


