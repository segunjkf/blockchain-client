# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o blockchain-client .

# Final stage
FROM alpine:latest

# Install ca-certificates for secure connections
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/blockchain-client .

# Expose port 8080
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Command to run the executable
CMD ["./blockchain-client"]