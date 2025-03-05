package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS" envDefault:"localhost:8080"`
	Key            string `env:"KEY" envDefault:""`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"3"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	RateLimit      int    `env:"RATE_LIMIT" envDefault:"3"`
}

type ServerConfig struct {
	Address         string `env:"ADDRESS" envDefault:"localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	DatabaseDsn     string `env:"DATABASE_DSN" envDefault:""`
	Key             string `env:"KEY" envDefault:""`
	StoreInterval   int    `env:"STORE_INTERVAL" envDefault:"300"`
	Restore         bool   `env:"RESTORE" envDefault:"true"`
}

func NewAgent(logger *zap.SugaredLogger) (*AgentConfig, error) {
	var cfg AgentConfig

	cfg, err := env.ParseAs[AgentConfig]()
	if err != nil {
		return &cfg, err
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

	logger.Infoln("Server Address:", cfg.ServerAddress)
	logger.Infoln("Report Interval:", cfg.ReportInterval)
	logger.Infoln("Poll Interval:", cfg.PollInterval)
	logger.Infoln("Secret Key:", cfg.Key)
	logger.Infoln("Rate limit:", cfg.RateLimit)
	return &cfg, nil
}

func NewServer(logger *zap.SugaredLogger) (*ServerConfig, error) {
	var cfg ServerConfig

	cfg, err := env.ParseAs[ServerConfig]()
	if err != nil {
		return &cfg, err
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

	logger.Infoln("Server Address:", cfg.Address)
	logger.Infoln("Store Interval:", cfg.StoreInterval)
	logger.Infoln("File Storage Path:", cfg.FileStoragePath)
	logger.Infoln("Restore:", cfg.Restore)
	logger.Infoln("Database Dsn", cfg.DatabaseDsn)
	logger.Infoln("Secret Key:", cfg.Key)

	return &cfg, nil
}
