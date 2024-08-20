package storage

import (
	"context"
	"database/sql"
	"errors"

	handler "github.com/Elvilius/go-musthave-metrics-tpl/internal/handlers"
	"github.com/Elvilius/go-musthave-metrics-tpl/internal/models"
)

type DBStorage struct {
	DB *sql.DB
}

func NewDBStorage(db *sql.DB) handler.Storager {
	return &DBStorage{DB: db}
}

func (db *DBStorage) Save(ctx context.Context, metric models.Metrics) error {
	if metric.MType == models.Counter {
		query := `
    		INSERT INTO metrics (id, m_type, delta) 
    		VALUES ($1, 'counter', $2)
    		ON CONFLICT (id, m_type) 
    		DO UPDATE SET delta = metrics.delta + EXCLUDED.delta;`
		_, err := db.DB.ExecContext(ctx, query, metric.ID, *metric.Delta)
		return err
	} else {
		query := `
			INSERT INTO metrics (id, m_type, value) 
			VALUES ($1, 'gauge', $2)
			ON CONFLICT (id, m_type) 
			DO UPDATE SET value = EXCLUDED.value;`
		_, err := db.DB.ExecContext(ctx, query, metric.ID, *metric.Value)
		return err
	}
}

func (db *DBStorage) Get(ctx context.Context, mType string, ID string) (models.Metrics, bool, error) {
	var metric models.Metrics
	row := db.DB.QueryRowContext(ctx, "SELECT id, m_type, value, delta from metrics WHERE m_type = $1 AND id = $2", mType, ID)

	err := row.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return metric, false, nil
		}
		return metric, false, err
	}
	return metric, true, nil
}

func (db *DBStorage) GetAll(ctx context.Context) ([]models.Metrics, error) {
	metrics := make([]models.Metrics, 0)

	row, err := db.DB.QueryContext(ctx, "SELECT id, m_type, value, delta from metrics")
	if err != nil {
		return metrics, err
	}
	if row.Err() != nil {
		return metrics, err
	}

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

func (db *DBStorage) Updates(ctx context.Context, metrics []models.Metrics) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		var query string
		var args []interface{}

		if metric.MType == models.Counter {
			query = `
			INSERT INTO metrics (id, m_type, delta) 
			VALUES ($1, 'counter', $2)
			ON CONFLICT (id, m_type) 
			DO UPDATE SET delta = metrics.delta + EXCLUDED.delta;`
			args = []interface{}{metric.ID, *metric.Delta}
		} else {
			query = `
			INSERT INTO metrics (id, m_type, value) 
			VALUES ($1, 'gauge', $2)
			ON CONFLICT (id, m_type) 
			DO UPDATE SET value = EXCLUDED.value;`
			args = []interface{}{metric.ID, *metric.Value}
		}

		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	errCommit := tx.Commit()
	if errCommit != nil {
		tx.Rollback()
		return errCommit
	} 
	return nil
}
