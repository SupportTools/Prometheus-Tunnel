# Use golang alpine image as the builder stage
FROM golang:1.22.4-alpine3.20 AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Fetch dependencies using go mod if your project uses Go modules
RUN go mod download

# Version and Git Commit build arguments
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

# Build the Go app with versioning information
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/supporttools/Prometheus-Tunnel/pkg/version.Version=$VERSION -X github.com/supporttools/Prometheus-Tunnel/pkg/version.GitCommit=$GIT_COMMIT -X github.com/supporttools/Prometheus-Tunnel/pkg/version.BuildTime=$BUILD_DATE" -o /bin/prometheus-tunnel
RUN chmod +x /bin/prometheus-tunnel

# Use ubuntu as the final image
FROM ubuntu:latest

# Install Common Dependencies
RUN apt-get update && \
    apt install -y \
    ca-certificates \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy our static executable.
COPY --from=builder /bin/prometheus-tunnel /bin/prometheus-tunnel

# Run the prometheus-tunnel binary.
ENTRYPOINT ["/bin/prometheus-tunnel"]
