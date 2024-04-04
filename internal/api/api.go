package api

import (
	"AWS_Trafiic_Monitor/internal/model"
	"AWS_Trafiic_Monitor/internal/utils"
	"encoding/json"
	"fmt"
	"time"
)

type API struct {
	InstanceName string
}

func NewAPI(name string) *API {
	return &API{name}
}

func (a *API) MetricDataMonth(name string) (*model.MetricData, error) {
	now := time.Now()
	// get first and last day of current month
	start, end := utils.BeginningOfDay(now), utils.EndOfDay(now)
	// calculate the period
	period := end.Sub(start)
	cmd := fmt.Sprintf(GetInstanceMetricData, a.InstanceName, name, start.Unix(), end.Unix(), period.Seconds())
	// execute the cmd
	data, err := utils.Execute(cmd)
	if err != nil {
		return nil, err
	}
	// unmarshal to model
	var md model.MetricData
	if err = json.Unmarshal(data, &md); err != nil {
		return nil, err
	}
	return &md, nil
}
