package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type AgentConfig struct {
	ServerAddress  string
	PollInterval   int
	ReportInterval int
	Key            string
	RateLimit      int
}

type ServerConfig struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DatabaseDsn     string
	Key             string
}

func getEnvOrDefaultString(envVar string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultBool(envVar string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(envVar); ok {
		if parsedValue, err := strconv.ParseBool(value); err == nil {
			return parsedValue
		}
	}
	return defaultValue
}

func getEnvOrDefaultInt(envVar string, defaultValue int) int {
	if value, ok := os.LookupEnv(envVar); ok {
		if parsedValue, err := strconv.Atoi(value); err == nil {
			return parsedValue
		}
	}
	return defaultValue
}

func NewAgent() *AgentConfig {
	cfg := &AgentConfig{
		ServerAddress:  getEnvOrDefaultString("ADDRESS", "localhost:8080"),
		PollInterval:   getEnvOrDefaultInt("POLL_INTERVAL", 3),
		ReportInterval: getEnvOrDefaultInt("REPORT_INTERVAL", 10),
		Key:            getEnvOrDefaultString("KEY", ""),
		RateLimit:      getEnvOrDefaultInt("RATE_LIMIT", 3),
	}

	pollInterval := flag.Int("p", cfg.PollInterval, "pollInterval")
	reportInterval := flag.Int("r", cfg.ReportInterval, "reportInterval")
	serverAddress := flag.String("a", cfg.ServerAddress, "server address")
	secretKey := flag.String("k", cfg.Key, "secret key")
	rateLimit := flag.Int("l", cfg.RateLimit, "rate limit")
	flag.Parse()

	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.ServerAddress = *serverAddress
	cfg.Key = *secretKey
	cfg.RateLimit = *rateLimit

	fmt.Println("Server Address:", cfg.ServerAddress)
	fmt.Println("Report Interval:", cfg.ReportInterval)
	fmt.Println("Poll Interval:", cfg.PollInterval)
	fmt.Println("Secret Key:", cfg.Key)
	fmt.Println("Rate limit:", cfg.RateLimit)
	return cfg
}

func NewServer() *ServerConfig {
	cfg := &ServerConfig{
		Address:         getEnvOrDefaultString("ADDRESS", "localhost:8080"),
		StoreInterval:   getEnvOrDefaultInt("STORE_INTERVAL", 300),
		FileStoragePath: getEnvOrDefaultString("FILE_STORAGE_PATH", "/tmp/metrics-db.json"),
		Restore:         getEnvOrDefaultBool("RESTORE", true),
		DatabaseDsn:     getEnvOrDefaultString("DATABASE_DSN", "postgres://root:qwerty@localhost:5432/metrics?sslmode=disable"),
		Key:             getEnvOrDefaultString("KEY", ""),
	}

	serverAddress := flag.String("a", cfg.Address, "server address")
	storeInterval := flag.Int("i", cfg.StoreInterval, "store interval")
	fileStoragePath := flag.String("f", cfg.FileStoragePath, "file storage path")
	restore := flag.Bool("r", cfg.Restore, "restore")
	databaseDsn := flag.String("d", cfg.DatabaseDsn, "database dsn")
	secretKey := flag.String("k", cfg.Key, "secret key")

	flag.Parse()

	cfg.Address = *serverAddress
	cfg.StoreInterval = *storeInterval
	cfg.FileStoragePath = *fileStoragePath
	cfg.Restore = *restore
	cfg.DatabaseDsn = *databaseDsn
	cfg.Key = *secretKey

	fmt.Println("Server Address:", cfg.Address)
	fmt.Println("Store Interval:", cfg.StoreInterval)
	fmt.Println("File Storage Path:", cfg.FileStoragePath)
	fmt.Println("Restore:", cfg.Restore)
	fmt.Println("Database Dsn", cfg.DatabaseDsn)
	fmt.Println("Secret Key:", cfg.Key)

	return cfg
}
