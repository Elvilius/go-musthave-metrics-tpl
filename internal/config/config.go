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
}

type ServerConfig struct {
	Address         string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
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

func GetAgentConfig() AgentConfig {
	cfg := AgentConfig{
		ServerAddress:  getEnvOrDefaultString("ADDRESS", "localhost:8080"),
		PollInterval:   getEnvOrDefaultInt("POLL_INTERVAL", 3),
		ReportInterval: getEnvOrDefaultInt("REPORT_INTERVAL", 10),
	}

	pollInterval := flag.Int("p", cfg.PollInterval, "pollInterval")
	reportInterval := flag.Int("r", cfg.ReportInterval, "reportInterval")
	serverAddress := flag.String("a", cfg.ServerAddress, "server address")
	flag.Parse()

	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.ServerAddress = *serverAddress

	fmt.Println("Server Address:", cfg.ServerAddress)
	fmt.Println("Report Interval:", cfg.ReportInterval)
	fmt.Println("Poll Interval:", cfg.PollInterval)

	return cfg
}

func GetServerConfig() ServerConfig {
	cfg := ServerConfig{
		Address:         getEnvOrDefaultString("ADDRESS", "localhost:8080"),
		StoreInterval:   getEnvOrDefaultInt("STORE_INTERVAL", 20),
		FileStoragePath: getEnvOrDefaultString("FILE_STORAGE_PATH", "/tmp/metrics-db.json"),
		Restore:         getEnvOrDefaultBool("RESTORE", true),
	}

	serverAddress := flag.String("a", cfg.Address, "server address")
	storeInterval := flag.Int("i", cfg.StoreInterval, "store interval")
	fileStoragePath := flag.String("f", cfg.FileStoragePath, "file storage path")
	restore := flag.Bool("r", cfg.Restore, "restore")

	flag.Parse()

	cfg.Address = *serverAddress
	cfg.StoreInterval = *storeInterval
	cfg.FileStoragePath = *fileStoragePath
	cfg.Restore = *restore

	fmt.Println("Server Address:", cfg.Address)
	fmt.Println("Store Interval:", cfg.StoreInterval)
	fmt.Println("File Storage Path:", cfg.FileStoragePath)
	fmt.Println("Restore:", cfg.Restore)

	return cfg
}
