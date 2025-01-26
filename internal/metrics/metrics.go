package metrics

import (
	"context"
	"strconv"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"go.uber.org/zap"
)

type Storager interface {
	Save(ctx context.Context, metric models.Metrics) error
	Get(ctx context.Context, mType, id string) (models.Metrics, bool, error)
	GetAll(ctx context.Context) ([]models.Metrics, error)
	Updates(ctx context.Context, metrics []models.Metrics) error
}

type Metrics struct {
	store  Storager
	logger *zap.SugaredLogger
}

func New(store Storager, logger *zap.SugaredLogger) *Metrics {
	return &Metrics{store: store, logger: logger}
}

func (m *Metrics) Add(ctx context.Context, metric models.Metrics, value string) error {
	if value != "" {
		if metric.MType == models.Counter {
			parseInt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			metric.Delta = &parseInt
		} else {
			parseFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			metric.Value = &parseFloat
		}
	}

	return m.store.Save(ctx, metric)
}

func (m *Metrics) GetOne(ctx context.Context, mType, id string) (*models.Metrics, error) {
	metric, ok, err := m.store.Get(ctx, mType, id)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &metric, nil
}

func (m *Metrics) GetAll(ctx context.Context) ([]models.Metrics, error) {
	return m.store.GetAll(ctx)
}

func (m *Metrics) Update(ctx context.Context, metrics []models.Metrics) error {
	return m.store.Updates(ctx, metrics)
}
