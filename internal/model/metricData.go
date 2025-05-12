package model

type MetricSum struct {
	Name string          `json:"metricName"`
	Data []MetricSumData `json:"metricData"`
}

type MetricSumData struct {
	Sum float64 `json:"sum"`
}
