FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the users service
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o users-service ./services/users/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/users-service .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8081

# Run the users service
CMD ["./users-service"] 