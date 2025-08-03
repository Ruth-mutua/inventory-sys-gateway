FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the gateway
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o gateway ./gateway/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/gateway .
COPY --from=builder /app/routes.yaml .

# Expose port
EXPOSE 8000

# Run the gateway
CMD ["./gateway"] 