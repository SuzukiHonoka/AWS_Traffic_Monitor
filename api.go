package main

var (
	MetricNames = []string{"NetworkIn", "NetworkOut"}
)

const (
	Cmd = "aws lightsail get-instance-metric-data --instance-name %s --metric-name %s --start-time %d --end-time %d --unit Bytes --statistics Sum --period %d"
)

type Data struct {
	MetricName string       `json:"metricName"`
	MetricData []metricData `json:"metricData"`
}

type metricData struct {
	Sum float64 `json:"sum"`
}
