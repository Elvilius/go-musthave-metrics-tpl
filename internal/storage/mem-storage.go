package storage

type MemStorage struct {
	metrics map[string]float64
}

func NewMemStorage() Storage {
	return &MemStorage{metrics: make(map[string]float64)}
}

func (r *MemStorage) Gauge(name string, value float64) {
	r.metrics[name] = value
}

func (r *MemStorage) Inc(name string) {
	_, ok := r.metrics[name]
	if !ok {
		r.metrics[name] = 1
	} else {
		r.metrics[name]++
	}
}
