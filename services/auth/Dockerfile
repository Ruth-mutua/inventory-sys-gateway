FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the auth service
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o auth-service ./services/auth/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/auth-service .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8083

# Run the auth service
CMD ["./auth-service"] 