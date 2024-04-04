package api

const (
	GetInstanceMetricData = "aws lightsail get-instance-metric-data --instance-name %s --metric-name %s --start-time %d --end-time %d --unit Bytes --statistics Sum --period %.0f"
)
