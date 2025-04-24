package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	Key            string `env:"KEY" envDefault:"" json:"key"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config         string `env:"CONFIG" json:"-"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"3" json:"poll_interval"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10" json:"report_interval"`
	RateLimit      int    `env:"RATE_LIMIT" envDefault:"3" json:"rate_limit"`
}

type ServerConfig struct {
	Address         string `env:"ADDRESS" envDefault:"localhost:8080" json:"address"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json" json:"store_file"`
	DatabaseDsn     string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`
	Key             string `env:"KEY" envDefault:"" json:"key"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config          string `env:"CONFIG" json:"-"`
	StoreInterval   int    `env:"STORE_INTERVAL" envDefault:"300" json:"store_interval"`
	Restore         bool   `env:"RESTORE" envDefault:"true" json:"restore"`
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
	cryptoKey := flag.String("c", cfg.CryptoKey, "crypto-key")
	config := flag.String("config", cfg.Config, "config")

	flag.Parse()

	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.ServerAddress = *serverAddress
	cfg.Key = *secretKey
	cfg.RateLimit = *rateLimit
	cfg.CryptoKey = *cryptoKey
	cfg.Config = *config

	if err := cfg.loadConfig(); err != nil {
		return nil, err
	}

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

	fmt.Println(cfg.Config, "CONFIG")
	serverAddress := flag.String("a", cfg.Address, "server address")
	storeInterval := flag.Int("i", cfg.StoreInterval, "store interval")
	fileStoragePath := flag.String("f", cfg.FileStoragePath, "file storage path")
	restore := flag.Bool("r", cfg.Restore, "restore")
	databaseDsn := flag.String("d", cfg.DatabaseDsn, "database dsn")
	secretKey := flag.String("k", cfg.Key, "secret key")
	cryptoKey := flag.String("c", cfg.CryptoKey, "crypto-key")
	config := flag.String("config", cfg.Config, "config")

	flag.Parse()

	cfg.Address = *serverAddress
	cfg.StoreInterval = *storeInterval
	cfg.FileStoragePath = *fileStoragePath
	cfg.Restore = *restore
	cfg.DatabaseDsn = *databaseDsn
	cfg.Key = *secretKey
	cfg.CryptoKey = *cryptoKey
	cfg.Config = *config

	// Load json config
	if err := cfg.loadConfig(); err != nil {
		return nil, err
	}

	logger.Infoln("Server Address:", cfg.Address)
	logger.Infoln("Store Interval:", cfg.StoreInterval)
	logger.Infoln("File Storage Path:", cfg.FileStoragePath)
	logger.Infoln("Restore:", cfg.Restore)
	logger.Infoln("Database Dsn", cfg.DatabaseDsn)
	logger.Infoln("Secret Key:", cfg.Key)
	logger.Infoln("Config:", cfg.Config)

	return &cfg, nil
}

func (s *ServerConfig) loadConfig() error {
	if s.Config == "" {
		return nil
	}

	data, err := loadConfig(s.Config)
	if err != nil {
		return err
	}

	loadCfg := &ServerConfig{}

	errUnmarshal := json.Unmarshal(data, loadCfg)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	if s.Address == "" {
		s.Address = loadCfg.Address
	}
	if s.CryptoKey == "" {
		s.CryptoKey = loadCfg.CryptoKey
	}
	if s.DatabaseDsn == "" {
		s.DatabaseDsn = loadCfg.DatabaseDsn
	}
	if s.FileStoragePath == "" {
		s.FileStoragePath = loadCfg.FileStoragePath
	}
	if s.Key == "" {
		s.Key = loadCfg.Key
	}
	if !s.Restore {
		s.Restore = loadCfg.Restore
	}
	if s.StoreInterval == 0 {
		s.StoreInterval = loadCfg.StoreInterval
	}

	return nil
}

func (a *AgentConfig) loadConfig() error {
	if a.Config == "" {
		return nil
	}

	data, err := loadConfig(a.Config)
	if err != nil {
		return err
	}

	loadCfg := &AgentConfig{}

	errUnmarshal := json.Unmarshal(data, loadCfg)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	if a.ServerAddress == "" {
		a.ServerAddress = loadCfg.ServerAddress
	}
	if a.CryptoKey == "" {
		a.CryptoKey = loadCfg.CryptoKey
	}

	if a.PollInterval == 0 {
		a.PollInterval = loadCfg.PollInterval
	}
	if a.Key == "" {
		a.Key = loadCfg.Key
	}
	if a.ReportInterval == 0 {
		a.ReportInterval = loadCfg.ReportInterval
	}
	if a.RateLimit == 0 {
		a.RateLimit = loadCfg.RateLimit
	}

	return nil
}

func loadConfig(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := make([]byte, 1024)

	var b bytes.Buffer

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		b.Write(buf[:n])
	}

	return b.Bytes(), nil

}
