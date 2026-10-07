# ---- Build Stage ----
# Use the official Go image as a builder.
# Specify the Go version to match your go.mod.
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies.
# This leverages Docker's layer caching for faster subsequent builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the Go application for Linux.
# CGO_ENABLED=0 creates a static binary which is ideal for containers.
# -ldflags="-w -s" strips debug information, making the binary smaller.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./src

# ---- Final Stage ----
# Use a minimal base image for a small and secure final image.
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the builder stage.
COPY --from=builder /app/server .
# Copy templates and static assets from the builder stage to ensure they are in the correct location.
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose the port the application runs on.
EXPOSE 8080

# Define the command to run the application when the container starts.
CMD ["./server"]