# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache \
  pkgconfig \
  vips-dev \
  gcc \
  musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -trimpath -o compactify .

FROM alpine:latest

RUN apk upgrade --no-cache && \
    apk add --no-cache vips

WORKDIR /workspace

COPY --from=builder /app/compactify /usr/local/bin/compactify

ENTRYPOINT ["compactify"]