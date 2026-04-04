package models

type LogsBatch struct {
	Logs []Log `json:"logs"`
}

type MetricsBatch struct {
	Metrics []Metric `json:"metrics"`
}