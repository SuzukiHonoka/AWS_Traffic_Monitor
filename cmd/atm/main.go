package main

import (
	"AWS_Trafiic_Monitor/internal/instance"
	"AWS_Trafiic_Monitor/internal/utils"
	"flag"
	"log"
	"time"
)

func init() {
	utils.CheckAwsCli()
}

func main() {
	// parse flags
	cPath := flag.String("c", "./config.json", "path to json config")
	loop := flag.Int("l", 0, "interval for loop")
	flag.Parse()
	path := *cPath
	// load config
	instances, err := instance.Load(path)
	utils.CheckError(err)
	instances.Check()
	if *loop > 0 {
		for {
			log.Println("Looping..")
			time.Sleep(time.Duration(*loop) * time.Second)
			instances.Check()
		}
	}
}
