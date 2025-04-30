# Use the official Golang image to create a build artifact.
FROM golang:1.23.5 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=linux go build -o main .

# Use a minimal base image
FROM debian:latest

# Install necessary packages for CGO
RUN apt-get update && apt-get install -y libc6

# Create a non-root user and group
RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built binary and .env file from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Change ownership of the files to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]