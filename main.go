package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	var instances []Instance
	cPath := flag.String("c", "", "path to json config")
	loop := flag.Int("l", 0, "interval for loop")
	flag.Parse()
	if len(*cPath) == 0 {
		fmt.Println("config not found")
		os.Exit(1)
	}
	if _, err := os.Stat(*cPath); os.IsNotExist(err) {
		fmt.Printf("Writing empty config file to %s", *cPath)
		b, err := json.Marshal([]Instance{{
			Name:    "",
			Limit:   Limit{},
			Command: []string{},
		}})
		checkError(err)
		err = ioutil.WriteFile(*cPath, b, os.ModePerm)
		checkError(err)
		fmt.Println("Config written")
		return
	}
	b, err := ioutil.ReadFile(*cPath)
	checkError(err)
	err = json.Unmarshal(b, &instances)
	checkError(err)
	now := time.Now()
	Check(instances, now)
	if *loop > 0 {
		for {
			fmt.Println("Looping..")
			time.Sleep(time.Duration(*loop) * time.Second)
			now = time.Now()
			Check(instances, now)
		}
	}
}

func BeginningOfDay(now time.Time) time.Time {
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

func EndOfDay(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
}

func BytesToUint(value float64, unit string) float32 {
	switch unit {
	case "GB":
		return float32(value / math.Pow(1024, 3))
	case "TB":
		return float32(value / math.Pow(1024, 4))
	default:
		panic(errors.New("unit not supported"))
	}
	return 0
}

func Check(units []Instance, time time.Time) {
	start, end := BeginningOfDay(time).Unix(), EndOfDay(time).Unix()
	fmt.Println("----------")
	for i, v := range units {
		var total int
		fmt.Printf("Instance Index: %d Name: %s\n", i, v.Name)
		for _, vv := range metricNames {
			var data Data
			checkError(json.Unmarshal([]byte(Exec(fmt.Sprintf(cmd, v.Name, vv, start, end, 2678400))), &data))
			used := int(BytesToUint(data.MetricData[0].Sum, v.Limit.Unit))
			total += used
			fmt.Printf("Metric: %s Used: %d %s\n", vv, used, v.Limit.Unit)
		}
		fmt.Printf("Total Used: %d %s\n", total, v.Limit.Unit)
		if total >= v.Limit.Value {
			more := total - v.Limit.Value
			fmt.Printf("Overflow: %d (%s)", more, strconv.Itoa(int(float32(more)/float32(v.Limit.Value)*100))+"%")
			fmt.Println("Trigger Fallback CMD..")
			for i, vvv := range v.Command {
				fmt.Printf("Command Index: %d CMD: [%s]\n", i, vvv)
				fmt.Printf("Excuted: [%s]\n", Exec(vvv))
			}
			fmt.Println("Exit Normally")
			os.Exit(0)
		} else {
			left := v.Limit.Value - total
			fmt.Printf("Traffic Left: %d %s (%s)\n", left, v.Limit.Unit, strconv.Itoa(int(float32(left)/float32(v.Limit.Value)*100))+"%")
		}
		fmt.Println("----------")
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Exec(cmd string) string {
	args := strings.Split(cmd, " ")
	result, err := exec.Command(args[0], args[1:]...).Output()
	checkError(err)
	return string(result)
}
