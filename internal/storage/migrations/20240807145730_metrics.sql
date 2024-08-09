-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics (
		    m_type VARCHAR(128) NOT NULL,
		    id VARCHAR(128) NOT NULL,
		    delta INT,
		    value DOUBLE PRECISION,
            PRIMARY KEY (id, m_type)
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
