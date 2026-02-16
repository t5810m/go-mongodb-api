package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	MongoURI string
	Port     string
	Timeout  time.Duration
}

var appConfig *Config

// isValidMongoURI validates MongoDB connection string format
func isValidMongoURI(uri string) bool {
	return strings.HasPrefix(uri, "mongodb://") || strings.HasPrefix(uri, "mongodb+srv://")
}

// isValidPort validates port number (1-65535)
func isValidPort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid port '%s': not a number", port)
	}
	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", portNum)
	}
	return nil
}

// Init loads configuration from environment variables
func Init() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}

	// Load and validate port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := isValidPort(port); err != nil {
		return nil, fmt.Errorf("invalid PORT configuration: %w", err)
	}

	// Load and validate timeout
	timeout := 30 * time.Second // default timeout
	if timeoutEnv := os.Getenv("TIMEOUT"); timeoutEnv != "" {
		if t, err := strconv.Atoi(timeoutEnv); err != nil {
			log.Printf("warning: invalid TIMEOUT value '%s', using default %v", timeoutEnv, timeout)
		} else if t <= 0 {
			log.Printf("warning: TIMEOUT must be positive, using default %v", timeout)
		} else {
			timeout = time.Duration(t) * time.Second
		}
	}

	// Load and validate MongoDB URI
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI not found in environment variables")
	}
	if !isValidMongoURI(mongoURI) {
		return nil, fmt.Errorf("invalid MONGO_URI: must start with 'mongodb://' or 'mongodb+srv://'")
	}

	appConfig = &Config{
		MongoURI: mongoURI,
		Port:     port,
		Timeout:  timeout,
	}

	log.Printf("configuration loaded: port=%s, timeout=%v", port, timeout)
	return appConfig, nil
}

// Get returns the global config instance, returns error if not initialized
func Get() (*Config, error) {
	if appConfig == nil {
		return nil, fmt.Errorf("configuration not initialized: call Init() first")
	}
	return appConfig, nil
}
