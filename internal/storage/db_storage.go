package storage

import (
	"context"
	"database/sql"
	"fmt"

	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
)

type DBStorage struct {
	Db *sql.DB
}

func NewDbStorage(db *sql.DB) handler.Storager {
	return &DBStorage{Db: db}
}

func (db *DBStorage) Save(ctx context.Context, metric models.Metrics) error {
	if metric.MType == models.Counter {
		query := `
    		INSERT INTO metrics (id, m_type, delta) 
    		VALUES ($1, 'counter', $2)
    		ON CONFLICT (id, m_type) 
    		DO UPDATE SET delta = metrics.delta + EXCLUDED.delta;`
		_, err := db.Db.ExecContext(ctx, query, metric.ID, *metric.Delta)
		if err != nil {
			fmt.Println(err)
		}
		return err
	} else {
		query := `
			INSERT INTO metrics (id, m_type, value) 
			VALUES ($1, 'gauge', $2)
			ON CONFLICT (id, m_type) 
			DO UPDATE SET value = EXCLUDED.value;`
		_, err := db.Db.ExecContext(ctx, query, metric.ID, *metric.Value)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
}

func (db *DBStorage) Get(ctx context.Context, mType string, ID string) (models.Metrics, bool, error) {
	var metric models.Metrics
	row := db.Db.QueryRowContext(ctx, "SELECT id, m_type, value, delta from metrics WHERE m_type = $1 AND id = $2", mType, ID)
	if row == nil {
		return metric, false, nil
	}
	err := row.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
	if err != nil {
		return metric, false, err
	}
	return metric, true, nil
}

func (db *DBStorage) GetAll(ctx context.Context) ([]models.Metrics, error) {
	row, _ := db.Db.QueryContext(ctx, "SELECT id, m_type, value, delta from metrics")
	metrics := make([]models.Metrics, 0)

	for row.Next() {
		var metric models.Metrics
		err := row.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
