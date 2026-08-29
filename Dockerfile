# Build
# FROM golang:1.25.3-alpine AS builder
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .

ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod download -x
RUN go build -o web-service ./cmd/server

# Run
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/web-service .
COPY templates ./templates
COPY static ./static

EXPOSE 8000
CMD ["./web-service"]