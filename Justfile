set quiet := true
set shell := ["bash", "-cu", "-o", "pipefail"]
set dotenv-load := true

NAME := 'dex'
IMAGE_NAME := 'dex_app'
REMOTE_DIR := 'dex'
DB_NAME := 'dex'

import? 'local.just'

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
    LOG_TYPE=none unbuffer go test -cover {{ package }} | gostack --test

cov package='./...': (gencov 'func' package)
cov-html package='./...': (gencov 'html' package)

[private]
gencov flag package: generate
    #!/bin/bash
    set -euo pipefail
    FILE=$(mktemp)
    export LOG_TYPE=none
    unbuffer go test {{ package }} -cover -coverprofile="$FILE" | gostack --test
    sed -i '/\.gen\.go/d' "$FILE"
    if [ "{{ flag }}" = func ]; then
        unbuffer go tool cover -{{ flag }}="$FILE" | gostack --test
    else
        go tool cover -{{ flag }}="$FILE"
    fi
    rm "$FILE"

lint:
    unbuffer go vet ./... | gostack
    unbuffer golangci-lint --color never run | gostack

fix:
    unbuffer golangci-lint --color never run --fix | gostack

check: test fmt lint

build: build-dex

build-dex: generate
    unbuffer go build -o bin/dex ./cmd/dex | gostack

buildc: generate
    docker-compose build

[private]
serve-dexd: generate
    unbuffer go run ./cmd/dexd | gostack

[private]
generate:
    go generate ./...
