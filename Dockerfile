FROM golang:1.25.5-alpine3.23 AS build
# https://hub.docker.com/_/golang

RUN apk add --no-cache \
    build-base \
    git

WORKDIR /build
COPY go.sum go.mod .
RUN go mod download

# Build.
COPY . .
RUN --mount=type=cache,id=dex,target=/root/.cache/go-build \
    go build -o dex -trimpath -ldflags '-s -w' ./cmd/dexd

# ——————————————————————————————————————————————————————————————————————————————————————————————————
FROM alpine:3.23.0
# https://hub.docker.com/_/alpine

RUN apk add --no-cache \
    bash \
    curl \
    mailcap \
    tzdata && \
    cp /usr/share/zoneinfo/Europe/London /etc/localtime && \
    echo 'Europe/London' > /etc/timezone
# mailcap for /etc/mime.types

RUN addgroup -g 1000 anon && \
    adduser -G anon -D -u 1000 anon

WORKDIR /app
COPY --from=build /build/dex /init

# This needs to be a script within the container because we need access to $PORT.
RUN echo 'curl -sSf http://localhost:$PORT/health' >> healthcheck && \
    chmod +x healthcheck

USER anon

CMD ["/init"]
