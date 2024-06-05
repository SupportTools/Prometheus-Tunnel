package config

import (
	"log"
	"os"
	"strconv"
)

// AppConfig structure for environment-based configurations.
type AppConfig struct {
	Debug       bool   `json:"debug"`
	MetricsPort int    `json:"metricsPort"`
	ServerIP    string `json:"ServerIP"`
	ServerPort  int    `json:"serverPort"`
}

// CFG is the global configuration instance populated by LoadConfiguration.
var CFG AppConfig

// LoadConfiguration loads the configuration from the environment variables.
func LoadConfiguration() {
	CFG.Debug = parseEnvBool("DEBUG", false)
	CFG.MetricsPort = parseEnvInt("METRICS_PORT", 9182)
	CFG.ServerIP = getEnvOrDefault("SERVER_IP", "")
	CFG.ServerPort = parseEnvInt("SERVER_PORT", 9182)
}

// getEnvOrDefault returns the value of the environment variable with the given key or the default value if the key is not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// parseEnvInt parses the environment variable with the given key and returns its integer representation or the default value if the key is
func parseEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing %s as int: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

// parseEnvBool parses the environment variable with the given key and returns its boolean representation or the default value if the key is not set.
func parseEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error parsing %s as bool: %v. Using default value: %t", key, err, defaultValue)
		return defaultValue
	}
	return boolValue
}
