package logger

import (
	"go.uber.org/zap"
)

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}
