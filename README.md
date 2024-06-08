# Prometheus Tunnel

[![Go Report Card](https://goreportcard.com/badge/github.com/supporttools/prometheus-tunnel)](https://goreportcard.com/report/github.com/supporttools/prometheus-tunnel)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/supporttools/prometheus-tunnel/CI)](https://github.com/supporttools/prometheus-tunnel/actions)
[![Docker Pulls](https://img.shields.io/docker/pulls/supporttools/prometheus-tunnel)](https://hub.docker.com/r/supporttools/prometheus-tunnel)
[![Helm Chart Version](https://img.shields.io/badge/helm%20chart-v0.1.0-blue)](https://charts.support.tools)
[![GitHub Release](https://img.shields.io/github/v/release/supporttools/prometheus-tunnel)](https://github.com/supporttools/prometheus-tunnel/releases/latest)
[![License](https://img.shields.io/github/license/supporttools/prometheus-tunnel)](LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FSupportTools%2FPrometheus-Tunnel.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FSupportTools%2FPrometheus-Tunnel?ref=badge_large&issueType=license)

Prometheus Tunnel is a reverse proxy server that forwards requests to a remote Prometheus exporter. This project includes functionality for metrics collection and health checks, making it easy to monitor and manage your Prometheus exporter.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Building from Source](#building-from-source)
  - [Using Docker](#using-docker)
  - [Using Helm](#using-helm)
- [Configuration](#configuration)
- [Prometheus Metrics](#prometheus-metrics)
- [Health Checks](#health-checks)
- [Prometheus Alert Rules](#prometheus-alert-rules)
- [Development](#development)
  - [Running Tests](#running-tests)
  - [Static Analysis Tools](#static-analysis-tools)
- [Contributing](#contributing)
- [License](#license)
- [Community and Support](#community-and-support)
- [Example Use Cases](#example-use-cases)

## [Features](#features)

- **Reverse Proxy**: Forwards incoming requests to a specified Prometheus exporter.
- **Metrics Collection**: Collects and exposes metrics for request count, duration, and response statuses.
- **Health Checks**: Provides endpoints for health and readiness checks.
- **Dockerized**: Easily deployable as a Docker container.
- **Helm Chart**: Package and deploy the application using Helm.

## [Quick Start](#quick-start)

Install the Prometheus Tunnel using Helm:

```bash
helm repo add supporttools https://charts.support.tools
helm upgrade --install prometheus-tunnel-server01 supporttools/prometheus-tunnel \
    --namespace monitoring \
    --create-namespace \
    --set settings.serverIP=192.168.0.3 \
    --set settings.serverPort=9182 \
    --set settings.name=server01
```

Note, please replace the following:

- `prometheus-tunnel-server01` with the name of the Helm release.
- `server_name` with the name of the remote Prometheus exporter (e.g., `server01`).
- `192.168.0.3` with the IP address of external Prometheus exporter.
- `9100` with the port of external Prometheus exporter.

## Installation

### Prerequisites

- Go 1.22 or higher
- Docker
- Kubernetes (optional, for deployment)
- Helm (optional, for deployment)

### Building from Source

1. Clone the repository:

    ```bash
    git clone https://github.com/supporttools/prometheus-tunnel.git
    cd prometheus-tunnel
    ```

2. Build the binary:

    ```bash
    go build -o prometheus-tunnel .
    ```

3. Run the application:

    ```bash
    ./prometheus-tunnel
    ```

### Using Docker

1. Build the Docker image:

    ```bash
    docker build -t supporttools/prometheus-tunnel:latest .
    ```

2. Run the Docker container:

    ```bash
    docker run -p 8080:8080 -e SERVER_IP=your-prometheus-exporter-ip -e SERVER_PORT=your-prometheus-exporter-port supporttools/prometheus-tunnel:latest
    ```

### Using Helm

1. Add the Helm repository:

    ```bash
    helm add repo supporttools https://charts.support.tools
    ```

2. Deploy using Helm:

    ```bash
    helm upgrade --install prometheus-tunnel supporttools/prometheus-tunnel \
        --namespace monitoring \
        --create-namespace \
        --values values.yaml
    ```

## Configuration

### Environment Variables

- `SERVER_IP`: IP address of the Prometheus exporter.
- `SERVER_PORT`: Port of the Prometheus exporter.
- `METRICS_PORT`: Port for exposing Prometheus metrics (default: 9182).
- `DEBUG`: Enable debug logging (default: false).

## Prometheus Metrics

The following metrics are exposed:

- `proxy_total_requests`: Total number of requests received.
- `proxy_request_duration_seconds`: Histogram of request durations.
- `proxy_response_status_total`: Count of responses by status code.

## Health Checks

- `/healthz`: Health check endpoint.
- `/readyz`: Readiness check endpoint.
- `/version`: Version information endpoint.

## Prometheus Alert Rules

Here is a sample `PrometheusRule` for monitoring the Prometheus Tunnel:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: prometheus-tunnel
  labels:
    prometheus: prometheus-tunnel
spec:
  groups:
  - name: prometheus-tunnel.rules
    rules:
    - alert: HighRequestRate
      expr: rate(proxy_total_requests[5m]) > 100
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High Request Rate"
        description: "The request rate has exceeded 100 requests per minute."

    - alert: SlowRequestDuration
      expr: histogram_quantile(0.99, rate(proxy_request_duration_seconds_bucket[5m])) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Slow Request Duration"
        description: "99th percentile request duration is greater than 1 second."

    - alert: HighErrorRate
      expr: rate(proxy_response_status_total{status=~"5.."}[5m]) > 10
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "High Error Rate"
        description: "The rate of 5xx errors has exceeded 10 errors per minute."

    - alert: HighLatency
      expr: rate(proxy_request_duration_seconds_sum[5m]) / rate(proxy_request_duration_seconds_count[5m]) > 0.5
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "High Request Latency"
        description: "The average request latency is greater than 0.5 seconds."

    - alert: HighRequestVolume
      expr: sum(rate(proxy_total_requests[5m])) by (job) > 1000
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High Request Volume"
        description: "The total request volume has exceeded 1000 requests per minute."
```

## Development

### Running Tests

To run the tests locally:

```bash
go test -v ./...
```

### Static Analysis Tools

To run static analysis tools:

```bash
golint ./...
staticcheck ./...
gosec ./...
```

## Contributing

Contributions are welcome! Please fork the repository and create a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Community and Support

For support and community interaction, you can join our Slack channel or open an issue on GitHub.

## Example Use Cases

- Prometheus servers managed by Argocd or FluxCD where editing the scape_configs to add external exporters is difficult.
- Air-gapped environments where the Prometheus exporter is not directly accessible from the cluster and a HTTP/Socks Proxy is required.
- Monitoring a Prometheus exporter that is behind a firewall and only accessible via a VPN.
