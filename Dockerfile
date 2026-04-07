FROM golang:1.26-alpine AS builder

RUN apk add --no-cache \
  pkgconfig \
  vips-dev \
  gcc \
  musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -trimpath -o compactify .

FROM alpine:latest

RUN apk add --no-cache vips

WORKDIR /workspace

COPY --from=builder /app/compactify /usr/local/bin/compactify

ENTRYPOINT ["compactify"]