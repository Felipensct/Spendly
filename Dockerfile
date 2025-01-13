# Use the official Go image
FROM golang:1.23.4 as builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the container
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o auth-service cmd/main.go

# Use a minimal image to run the service
FROM alpine:3.18

WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/auth-service .

# Expose the gRPC port
EXPOSE 50051

# Command to run the binary
CMD ["./auth-service"]
