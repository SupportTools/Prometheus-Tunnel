package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
	proxyHandler, err := proxy.NewProxyHandler(config.CFG.ServerIP, config.CFG.ServerPort, config.CFG.Debug)
	if err != nil {
		log.Fatalf("Failed to create proxy handler: %v", err)
	}

	go func() {
		log.Printf("Starting proxy server on port %d\n", config.CFG.ServerPort)
		srv := &http.Server{
			Addr:         ":" + strconv.Itoa(config.CFG.ServerPort),
			Handler:      http.HandlerFunc(proxyHandler),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigs
	log.Infoln("Received shutdown signal, shutting down proxy server...")
	done <- true
	log.Infoln("Prometheus-Tunnel shut down gracefully")
}
