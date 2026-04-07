FROM golang:1.26-bookworm AS builder

RUN apt-get update && apt-get install -y \
  pkg-config \
  libvips-dev \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -trimpath -o compactify .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
  libvips \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY --from=builder /app/compactify /usr/local/bin/compactify

ENTRYPOINT ["compactify"]