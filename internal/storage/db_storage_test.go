package storage

import (
	"context"
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/mocks"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDBStorage_SaveCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorager(ctrl)

	ctx := context.Background()

	delta := int64(10)
	testMetric := models.Metrics{
		ID:    "test_metric",
		MType: models.Counter,
		Delta: &delta,
	}

	db.EXPECT().Save(ctx, testMetric).Return(nil)

	err := db.Save(ctx, testMetric)

	assert.NoError(t, err)
}

func TestDBStorage_SaveGauge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorager(ctrl)

	ctx := context.Background()

	value := 10.11
	testMetric := models.Metrics{
		ID:    "test_metric",
		MType: models.Counter,
		Value: &value,
	}

	db.EXPECT().Save(ctx, testMetric).Return(nil)

	err := db.Save(ctx, testMetric)

	assert.NoError(t, err)
}

func TestDBStorage_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorager(ctrl)

	ctx := context.Background()

	delta := int64(10)

	expectedMetric := models.Metrics{
		ID:    "test_metric",
		MType: models.Counter,
		Delta: &delta,
	}

	db.EXPECT().Get(ctx, "test_metric", models.Counter).Return(expectedMetric, true, nil)

	result, found, err := db.Get(ctx, "test_metric", models.Counter)

	assert.NoError(t, err)

	assert.True(t, found)

	assert.Equal(t, expectedMetric, result)
}

func TestDBStorage_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockStorager(ctrl)

	ctx := context.Background()

	delta1 := int64(10)
	delta2 := int64(20)

	expectedMetrics := []models.Metrics{
		{
			ID:    "metric1",
			MType: models.Counter,
			Delta: &delta1,
		},
		{
			ID:    "metric2",
			MType: models.Counter,
			Delta: &delta2,
		},
	}

	db.EXPECT().GetAll(ctx).Return(expectedMetrics, nil)

	metrics, err := db.GetAll(ctx)

	assert.NoError(t, err)

	assert.Equal(t, expectedMetrics, metrics)
}
