package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

// NewProxyHandler creates a new proxy handler function
func NewProxyHandler(serverIp string, serverPort int, debug bool) (func(w http.ResponseWriter, r *http.Request), error) {
	serverURL := "http://" + serverIp + ":" + strconv.Itoa(serverPort)
	proxyURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)

	return func(w http.ResponseWriter, r *http.Request) {
		if debug {
			logRequest(r)
		}
		r.Host = proxyURL.Host
		proxy.ServeHTTP(w, r)
	}, nil
}

func logRequest(r *http.Request) {
	log.Printf("Received request: %s %s %s\n", r.Method, r.URL, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("Remote Address: %s\n", r.RemoteAddr)
	log.Printf("Request URI: %s\n", r.RequestURI)
	log.Printf("Headers:\n")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s\n", name, value)
		}
	}
	log.Printf("Timestamp: %s\n", time.Now().Format(time.RFC3339))
}
