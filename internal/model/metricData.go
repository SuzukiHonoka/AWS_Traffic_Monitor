package model

type MetricData struct {
	Name string       `json:"metricName"`
	Data []metricData `json:"metricData"`
}

type metricData struct {
	Sum float64 `json:"sum"`
}
