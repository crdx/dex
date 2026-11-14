set quiet := true
set shell := ["bash", "-cu", "-o", "pipefail"]
set dotenv-load := true

NAME := 'dex'
IMAGE_NAME := 'dex_app'
REMOTE_DIR := 'dex'
DB_NAME := 'dex'

import? 'local.just'

mod release

export CGO_ENABLED := '0'

mod make

[private]
help:
    just --list --unsorted --list-submodules

dev:
    hivemind

db:
    mariadb {{ DB_NAME }}

drop:
    echo 'drop database if exists {{ DB_NAME }}' | mariadb

fmt:
    go fmt ./...

test package='./...': generate
    LOG_TYPE=none go test -cover {{ package }}

cov package='./...': (gencov 'func' package)
cov-html package='./...': (gencov 'html' package)

[private]
gencov flag package: generate
    #!/bin/bash
    set -euo pipefail
    FILE=$(mktemp)
    export LOG_TYPE=none
    go test {{ package }} -cover -coverprofile="$FILE"
    sed -i '/\.gen\.go/d' "$FILE"
    if [ "{{ flag }}" = func ]; then
        go tool cover -{{ flag }}="$FILE"
    else
        go tool cover -{{ flag }}="$FILE"
    fi
    rm "$FILE"

lint:
    go vet ./...
    golangci-lint run --color never

fix:
    golangci-lint run --color never --fix

check: test fmt lint

build: build-dex

build-dex: generate
    go build -o bin/dex ./cmd/dex

buildc: generate
    docker-compose build

[private]
serve-dexd: generate
    go run ./cmd/dexd

[private]
generate:
    go generate ./...
