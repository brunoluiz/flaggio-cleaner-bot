FROM alpine:3.12
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
COPY flaggio-cleaner-bot /
ENTRYPOINT ["/flaggio-cleaner-bot"]
