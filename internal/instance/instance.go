package instance

import (
	"AWS_Trafiic_Monitor/internal/api"
	"AWS_Trafiic_Monitor/internal/utils"
	"log"
	"os"
)

var MetricNames = []string{"NetworkIn", "NetworkOut"}

type Instance struct {
	Name    string
	Limit   Limit
	Command []string
}

func (x *Instance) Check() {
	log.Println("----------")
	log.Printf("Instance Name: %s", x.Name)

	newAPI := api.NewAPI(x.Name)
	total := 0
	for _, name := range MetricNames {
		md, err := newAPI.MetricDataMonth(name)
		if err != nil {
			log.Fatalf("Get data of metric: %s failed, err=%s", name, err)
		}

		if len(md.Data) == 0 {
			log.Printf("No data of metric: %s", name)
			continue
		}

		used := int(x.Limit.Unit.BytesToUnit(md.Data[0].Sum))
		total += used
		log.Printf("Metric: %s Used: %d %s\n", name, used, x.Limit.Unit)
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
		log.Println("Trigger Fallback commands..")
		x.ExecuteCommand()
		log.Println("Exiting..")
		os.Exit(0)
	}
	log.Println("----------")
}

func (x *Instance) ExecuteCommand() {
	for i, cmd := range x.Command {
		log.Printf("Command Index: %d CMD: [%s]", i, cmd)
		b, err := utils.Execute(cmd)
		if err != nil {
			log.Printf("Execute command: %s failed, err=%s", cmd, err)
			continue
		}
		log.Printf("Commaned: %s excuted: %s", cmd, string(b))
	}
}
