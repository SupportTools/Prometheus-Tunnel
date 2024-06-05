package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/supporttools/prometheus-tunnel/pkg/config"
	"github.com/supporttools/prometheus-tunnel/pkg/logging"
	"github.com/supporttools/prometheus-tunnel/pkg/metrics"
	"github.com/supporttools/prometheus-tunnel/pkg/proxy"
)

func main() {
	flag.Parse()
	config.LoadConfiguration()
	log := logging.SetupLogging(config.CFG.Debug)
	log.Debug("Debug logging enabled")

	log.Info("Starting Prometheus-Tunnel...")

	// Start Prometheus metrics server
	log.Infoln("Starting metrics server")
	go metrics.StartMetricsServer()

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	// Create and start proxy server
	proxyHandler, err := proxy.NewProxyHandler(config.CFG.ServerIp, config.CFG.ServerPort, config.CFG.Debug)
	if err != nil {
		log.Fatalf("Failed to create proxy handler: %v", err)
	}

	go func() {
		log.Printf("Starting proxy server on port %d\n", config.CFG.ServerPort)
		if err := http.ListenAndServe(":"+strconv.Itoa(config.CFG.ServerPort), http.HandlerFunc(proxyHandler)); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigs
	log.Infoln("Received shutdown signal, shutting down proxy server...")
	done <- true
	log.Infoln("Prometheus-Tunnel shut down gracefully")
}
