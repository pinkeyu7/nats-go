# Stage 1: Build stage
FROM golang:1.25-alpine AS build

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Add gcc
RUN apk add --no-cache gcc musl-dev

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./agent

# Stage 2: Final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main .

# Set the entrypoint command
ENTRYPOINT ["/app/main"]