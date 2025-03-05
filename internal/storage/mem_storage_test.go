package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/metrics"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/stretchr/testify/assert"
)

func saveMetric(store metrics.Storager, metrics []models.Metrics) {
	for _, m := range metrics {
		fmt.Println(store.Save(context.TODO(), m))
	}
}

func TestMemStorage_Save(t *testing.T) {
	ctx := context.Background()
	store := NewMemStorage()

	value := 123.2

	type args struct {
		ctx    context.Context
		metric models.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save success metric Gauge",
			args: args{
				ctx:    ctx,
				metric: models.Metrics{ID: models.MetricAlloc, MType: models.Gauge, Value: &value},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Save(ctx, tt.args.metric)
			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
				return
			}

			_, ok, _ := store.Get(ctx, tt.args.metric.MType, tt.args.metric.ID)
			assert.True(t, ok)

		})
	}
}

func TestMemStorage_Get(t *testing.T) {

	ctx := context.Background()
	store := NewMemStorage()
	value := 10.1

	saveMetric(store, []models.Metrics{
		{ID: models.MetricAlloc, MType: models.Gauge, Value: &value},
		{ID: models.MetricBuckHashSys, MType: models.Gauge, Value: &value},
		{ID: models.MetricFreeMemory, MType: models.Gauge, Value: &value},
		{ID: models.MetricCPUUtilization1, MType: models.Gauge, Value: &value},
	})

	type args struct {
		mType string
		id    string
	}

	tests := []struct {
		name string
		args args
		ok   bool
	}{
		{
			name: "Get success metric MetricAlloc",
			args: args{
				mType: models.Gauge,
				id:    models.MetricAlloc,
			},
			ok: true,
		},
		{
			name: "Get success metric MetricBuckHashSys",
			args: args{
				mType: models.Gauge,
				id:    models.MetricBuckHashSys,
			},
			ok: true,
		},
		{
			name: "Get success metric MetricFreeMemory",
			args: args{
				mType: models.Gauge,
				id:    models.MetricFreeMemory,
			},
			ok: true,
		},
		{
			name: "Get success metric MetricCPUUtilization1",
			args: args{
				mType: models.Gauge,
				id:    models.MetricCPUUtilization1,
			},
			ok: true,
		},

		{
			name: "Get fail metric unknown type",
			args: args{
				mType: "Test",
				id:    models.MetricCPUUtilization1,
			},
			ok: false,
		},
		{
			name: "Get fail metric unknown ID",
			args: args{
				mType: "Test",
				id:    "iiii",
			},
			ok: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok, _ := store.Get(ctx, tt.args.mType, tt.args.id)
			assert.Equal(t, tt.ok, ok)

		})
	}
}
