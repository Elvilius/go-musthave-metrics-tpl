package logger

import (
	"go.uber.org/zap"
)

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	defer func() {
		if errSyncLogger := logger.Sync(); errSyncLogger != nil {
			logger.Sugar().Errorln(errSyncLogger)
		}
	}()

	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}
