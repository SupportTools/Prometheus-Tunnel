# Start from the official Go image
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o Prometheus-Tunnel .

# Start a new stage from scratch
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/Prometheus-Tunnel /Prometheus-Tunnel

# Expose port 9000 to the outside world
EXPOSE 9000

# Command to run the executable
CMD ["/Prometheus-Tunnel"]
