# syntax=docker/dockerfile:1.7-labs

# UI
FROM --platform=$BUILDPLATFORM oven/bun:1-alpine AS js
WORKDIR /ui
COPY ui/bun.lock ui/package.json ./
RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install --frozen-lockfile
COPY ui ./
RUN bun run build

FROM --platform=$BUILDPLATFORM golang:alpine AS goenv
ARG TARGETOS
ARG TARGETARCH
RUN printf 'export GOOS=%s\nexport GOARCH=%s\n' "$TARGETOS" "$TARGETARCH" > /goenv

FROM --platform=$BUILDPLATFORM golang:alpine AS gomod
RUN apk add --no-cache git
ARG GOFLAGS
ARG GOEXPERIMENT
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM --platform=$BUILDPLATFORM gomod AS build-restic
ARG RESTIC_VERSION=0.19.0
ARG GOFLAGS
ARG GOEXPERIMENT
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,from=goenv,source=/goenv,target=/goenv \
    . /goenv && \
    go mod download "github.com/restic/restic@v${RESTIC_VERSION}" && \
    cd "$(go env GOMODCACHE)/github.com/restic/restic@v${RESTIC_VERSION}" && \
    CGO_ENABLED=0 go build -trimpath -o /out/restic ./cmd/restic

FROM --platform=$BUILDPLATFORM gomod AS build-server
ARG GOFLAGS
ARG GOEXPERIMENT
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,target=/src \
    --mount=type=bind,from=js,source=/ui/dist,target=/ui-dist \
    --mount=type=bind,from=goenv,source=/goenv,target=/goenv \
    . /goenv && \
    cp -a /src/. /work && \
    rm -rf /work/ui/dist && \
    mkdir -p /work/ui && \
    cp -a /ui-dist /work/ui/dist && \
    cd /work && \
    VERSION=$(git describe --tags --dirty --always 2>/dev/null || echo unknown) && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -extldflags '-static' -X main.Version=${VERSION}" -o /out/packup-server .

FROM --platform=$BUILDPLATFORM gomod AS build-agent
ARG GOFLAGS
ARG GOEXPERIMENT
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,target=/src \
    --mount=type=bind,from=goenv,source=/goenv,target=/goenv \
    . /goenv && \
    cd /src && \
    VERSION=$(git describe --tags --dirty --always 2>/dev/null || echo unknown) && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -extldflags '-static' -X main.Version=${VERSION}" -o /out/packup-agent ./agent

FROM alpine:latest AS base
RUN apk add --no-cache coreutils ca-certificates
COPY --from=build-restic /out/restic /usr/local/bin/restic
WORKDIR /backups

# agent
FROM base AS agent
RUN apk add --no-cache sqlite postgresql18-client mariadb-client rclone
COPY --from=build-agent /out/packup-agent /usr/local/bin/packup-agent
COPY ./mods /mods
RUN chmod +x /mods/*
ENTRYPOINT ["packup-agent"]

# server
FROM base AS server
COPY --from=build-server /out/packup-server /usr/local/bin/packup-server
ENTRYPOINT ["packup-server"]
CMD ["--config=/etc/packup/packup.yaml"]

# explicit default target
FROM server
