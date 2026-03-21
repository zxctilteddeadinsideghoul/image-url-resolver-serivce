# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd .

RUN go build -o image-resolver ./cmd

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/image-resolver .

EXPOSE 8080

CMD ["./image-resolver"]