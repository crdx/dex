FROM golang:1.26.3-trixie AS build
# https://hub.docker.com/_/golang

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        media-types \
        tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# Build.
COPY . .
RUN --mount=type=cache,id=dex,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o dex -trimpath -ldflags '-s -w' ./cmd/dexd

# ——————————————————————————————————————————————————————————————————————————————————————————————————
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo/Europe/London /etc/localtime
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/mime.types /etc/mime.types
COPY --from=build /build/dex /init

WORKDIR /app
USER 1000:1000

CMD ["/init"]
