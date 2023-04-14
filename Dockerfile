# Dockerfile

# Use the official Golang base image
FROM golang:1.20.3 as builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Start a new stage
FROM golang:1.20.3

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the static folder from the builder stage
COPY --from=builder /app/static /app/static

# Expose the port for the chat service
EXPOSE 8080

# Set the entry point for the container
ENTRYPOINT ["/app/main"]