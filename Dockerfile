# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1001 appuser && \
    adduser -u 1001 -G appuser -s /bin/sh -D appuser
WORKDIR /app
COPY --from=builder /app/main .
RUN chown appuser:appuser /app/main
USER appuser

EXPOSE 8080
CMD ["./main"] 