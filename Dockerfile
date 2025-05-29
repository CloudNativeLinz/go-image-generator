# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o go-image-generator ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/go-image-generator ./
COPY assets ./assets
COPY _data ./_data
COPY assets/fonts ./assets/fonts
COPY assets/backgrounds ./assets/backgrounds
COPY assets/overlays ./assets/overlays
COPY assets/templates ./assets/templates
COPY README.md ./README.md

ENTRYPOINT ["./go-image-generator"]
CMD ["--help"]
