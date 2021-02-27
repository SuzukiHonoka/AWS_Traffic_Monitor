package main

type Instance struct {
	Name string
	Limit Limit
	Command []string
}

type Limit struct {
	Unit string
	Value int
}

type Data struct {
	MetricName string     `json:"metricName"`
	MetricData []metricData `json:"metricData"`
}

type metricData struct {
	Sum float64 `json:"sum"`
}

var (
	metricNames = []string { "NetworkIn", "NetworkOut" }
)

const (
	cmd = "aws lightsail get-instance-metric-data --instance-name %s --metric-name %s --start-time %d --end-time %d --unit Bytes --statistics Sum --period %d"
)
