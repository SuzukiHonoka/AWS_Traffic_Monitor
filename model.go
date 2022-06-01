package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Instance struct {
	Name    string
	Limit   Limit
	Command []string
}

type Instances []Instance

func (x *Instance) Check() {
	now := time.Now()
	// get first and last day of this month
	start, end := BeginningOfDay(now).Unix(), EndOfDay(now).Unix()
	log.Println("----------")
	var total int
	log.Printf("Instance Name: %s", x.Name)
	for _, vv := range MetricNames {
		var data Data
		err := json.Unmarshal(Exec(fmt.Sprintf(Cmd, x.Name, vv, start, end, 2678400)), &data)
		checkError(err)
		if data.MetricData == nil {
			continue
		}
		used := int(BytesToUnit(x.Limit.Unit, data.MetricData[0].Sum))
		total += used
		log.Printf("Metric: %s Used: %d %s\n", vv, used, x.Limit.Unit)
	}
	log.Printf("Total Used: %d %s", total, x.Limit.Unit)
	if total >= x.Limit.Value {
		more := total - x.Limit.Value
		log.Printf("Overflow: %d (%.2f %%)", more, float32(more)/float32(x.Limit.Value)*100)
		log.Println("Trigger Fallback CMD..")
		for i, cmd := range x.Command {
			log.Printf("Command Index: %d CMD: [%s]", i, cmd)
			log.Printf("Excuted: [%s]", Exec(cmd))
		}
		log.Println("Exit Normally")
		os.Exit(0)
	} else {
		left := x.Limit.Value - total
		log.Printf("Traffic Left: %d %s (%.2f %%)", left, x.Limit.Unit, float32(left)/float32(x.Limit.Value)*100)
	}
	log.Println("----------")
}

func (s Instances) Check() {
	for _, v := range s {
		v.Check()
	}
}
