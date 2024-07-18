package logger

import "go.uber.org/zap"

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}
