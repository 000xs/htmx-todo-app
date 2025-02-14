# Use Golang image to build the Go app
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod ./
COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o myapp .

# Start a new stage from scratch to keep the image size smaller
FROM debian:bullseye-slim

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary from the builder stage
COPY --from=builder /app/myapp .

# Expose port 8080 (or whichever port your app uses)
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
