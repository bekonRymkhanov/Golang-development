# Start from the official Golang image.
FROM golang:1.22.1 as builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Optionally copy any other necessary files (e.g., migrations).
# If you have migrations or other data, uncomment and modify the next line.
# COPY path_to_migrations_folder ./path_in_container

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the application, disable CGO to create a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o project ./cmd/api

# Use a smaller image to run the app.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/project .
# If you copied migrations or other data, make sure to copy them to the final image.
COPY --from=builder /app/migrations ./migrations

# Command to run the executable.
CMD ["./project"]