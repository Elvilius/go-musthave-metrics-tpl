package repo

import "fmt"

type Repo struct {
	metrics map[string]float64
}

func NewRepo() *Repo {
	return &Repo{metrics: make(map[string]float64)}
}

func (r *Repo) Gauge(name string, value float64) {
	r.metrics[name] = value
}

func (r *Repo) Inc(name string) {
	_, ok := r.metrics[name]
	if !ok {
		r.metrics[name] = 1
	} else {
		r.metrics[name]++
	}
}


func (r *Repo) Print() {
	fmt.Println(123123123123)
	fmt.Println(r.metrics)
}


