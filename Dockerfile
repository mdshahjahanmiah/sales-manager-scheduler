# Stage 1: Build the Go application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application binary
RUN CGO_ENABLED=0 go build -o main ./cmd/main.go

# Stage 2: Final runtime image
FROM debian:bookworm-slim

# Set the working directory
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port that the application listens on
EXPOSE 3000

# Command to run the Go application
CMD ["./main"]