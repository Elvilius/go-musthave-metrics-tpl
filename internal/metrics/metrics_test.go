package metrics

import (
	"context"
	"fmt"
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/mocks"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_Add(t *testing.T) {
	ctx := context.Background()
	store := mocks.NewMockMemStorage()
	m := New(store, nil)

	type args struct {
		ctx    context.Context
		metric models.Metrics
		value  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add success metric Gauge",
			args: args{
				ctx:    ctx,
				metric: models.Metrics{ID: models.MetricAlloc, MType: models.Gauge},
				value:  "123.234",
			},
			wantErr: false,
		},
		{
			name: "Add success metric Counter",
			args: args{
				ctx:    ctx,
				metric: models.Metrics{ID: models.MetricPollCount, MType: models.Counter},
				value:  "100",
			},
			wantErr: false,
		},
		{
			name: "Add fail unknown type",
			args: args{
				ctx:    ctx,
				metric: models.Metrics{ID: models.MetricAlloc, MType: "TEST"},
				value:  "233.3",
			},
			wantErr: true,
		},
		{
			name: "Add fail invalid value",
			args: args{
				ctx:    ctx,
				metric: models.Metrics{ID: models.MetricAlloc, MType: models.Counter},
				value:  "tesdfsdf,",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.Add(tt.args.ctx, tt.args.metric, tt.args.value)
			fmt.Println(err)
			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
			}
		})
	}
}

func TestMetrics_Update(t *testing.T) {
	value := 123.2
	ctx := context.Background()
	store := mocks.NewMockMemStorage()
	m := New(store, nil)

	type args struct {
		ctx    context.Context
		metric []models.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add success metric Gauge",
			args: args{
				ctx:    ctx,
				metric: []models.Metrics{{ID: models.MetricAlloc, MType: models.Gauge, Value: &value}, {ID: models.MetricAlloc, MType: models.Gauge, Value: &value}},
			},
			wantErr: false,
		},

		{
			name: "Add fail invalid value",
			args: args{
				ctx:    ctx,
				metric: []models.Metrics{{ID: models.MetricAlloc, MType: "REDSD", }},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.Update(tt.args.ctx, tt.args.metric)
			fmt.Println(err)
			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
			}
		})
	}
}
