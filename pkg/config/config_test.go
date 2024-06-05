package config

import (
	"os"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	os.Setenv("DEBUG", "true")
	os.Setenv("METRICS_PORT", "9182")
	os.Setenv("SERVER_IP", "192.168.3.1")
	os.Setenv("SERVER_PORT", "9182")

	LoadConfiguration()

	if !CFG.Debug {
		t.Errorf("Expected Debug to be true, got %v", CFG.Debug)
	}
	if CFG.MetricsPort != 8080 {
		t.Errorf("Expected MetricsPort to be 8080, got %d", CFG.MetricsPort)
	}
	if CFG.ServerIp != "192.168.3.1" {
		t.Errorf("Expected ServerIp to be '192.168.3.1', got %v", CFG.ServerIp)
	}
	if CFG.ServerPort != 9182 {
		t.Errorf("Expected ServerPort to be 9182, got %d", CFG.ServerPort)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	key := "TEST_KEY"
	defaultValue := "default"

	value := getEnvOrDefault(key, defaultValue)
	if value != defaultValue {
		t.Errorf("Expected '%s', got '%s'", defaultValue, value)
	}

	expectedValue := "test-value"
	os.Setenv(key, expectedValue)
	value = getEnvOrDefault(key, defaultValue)
	if value != expectedValue {
		t.Errorf("Expected '%s', got '%s'", expectedValue, value)
	}
}
