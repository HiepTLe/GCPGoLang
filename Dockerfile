FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN make build

# Create a minimal image for running the application
FROM alpine:3.19

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/gcpgolang /app/gcpgolang

# Set the binary as executable
RUN chmod +x /app/gcpgolang

# Expose port (if needed)
# EXPOSE 8080

# Set environment variables
ENV GO_ENV=production

# Command to run the executable
ENTRYPOINT ["/app/gcpgolang"]

# Default arguments
CMD ["--help"] 