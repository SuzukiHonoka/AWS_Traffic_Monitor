package instance

import (
	"github.com/SuzukiHonoka/AWS-Traffic-Monitor/internal/api"
	"github.com/SuzukiHonoka/AWS-Traffic-Monitor/internal/utils"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"log"
	"time"
)

var MetricNameList = []types.InstanceMetricName{types.InstanceMetricNameNetworkIn, types.InstanceMetricNameNetworkOut}

type Instance struct {
	Name        string
	Limit       Limit
	CommandList []string
	api         *api.API
}

func (x *Instance) Check() {
	log.Println("----------")
	log.Printf("Instance Name: %s", x.Name)

	if x.api == nil {
		x.api = api.NewAPI(x.Name)
	}

	total := 0
	for _, name := range MetricNameList {
		md, err := x.api.MetricDataMonth(name)
		if err != nil {
			log.Printf("Get data of metric: %s failed, err=%s", name, err)
			continue
		}

		if len(md.Data) == 0 {
			log.Printf("No data of metric: %s", name)
			continue
		}

		used := int(x.Limit.Unit.FromBytes(md.Data[0].Sum))
		total += used
		log.Printf("Metric: %s Used: %d %s", name, used, x.Limit.Unit)
	}
	log.Printf("Total Used: %d %s", total, x.Limit.Unit)

	if total < x.Limit.Value {
		left := x.Limit.Value - total
		log.Printf("Traffic Left: %d %s (%.2f %%)", left, x.Limit.Unit, float64(left)/float64(x.Limit.Value)*100)
	} else {
		overflow := total - x.Limit.Value
		if overflow > 0 {
			log.Printf("Overflow: %d %s (%.2f %%)", overflow, x.Limit.Unit, float64(total)/float64(x.Limit.Value)*100)
		}

		log.Printf("Executing commands, count=%d", len(x.CommandList))
		now := time.Now()
		x.ExecuteCommand()
		log.Printf("Command executed, cost=%s", time.Since(now))
	}
	log.Println("----------")
}

func (x *Instance) ExecuteCommand() {
	for i, cmd := range x.CommandList {
		if cmd == "shutdown" {
			err := x.api.Shutdown(true)
			if err != nil {
				log.Printf("Shutdown instance: %s failed, err=%s", x.Name, err)
			} else {
				log.Printf("Shutdown instance: %s success", x.Name)
			}
			continue
		}

		log.Printf("Command index: %d -> [ %s ]", i, cmd)
		b, err := utils.Execute(cmd)
		if err != nil {
			log.Printf("Execute command: %s failed, err=%s", cmd, err)
			continue
		}
		log.Printf("Command output: %s", string(b))
	}
}
