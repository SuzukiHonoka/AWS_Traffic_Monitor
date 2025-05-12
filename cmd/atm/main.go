package main

import (
	"AWS_Trafiic_Monitor/internal/instance"
	"AWS_Trafiic_Monitor/internal/utils"
	"flag"
	"log"
	"time"
)

var (
	configPath   = flag.String("c", "./config.json", "path to json config")
	loopInterval = flag.Int("l", 0, "interval for loop")
)

func init() {
	// check if aws cli is installed
	utils.CheckAwsCli()

	// parse flags
	flag.Parse()
}

func main() {
	// parse flags
	path := *configPath

	// load config
	instances, err := instance.Load(path)
	if err != nil {
		log.Fatalf("failed to load config, err=%v", err)
	}

	instances.Check()
	if interval := *loopInterval; interval > 0 {
		for {
			log.Println("Looping..")
			time.Sleep(time.Duration(interval) * time.Second)
			instances.Check()
		}
	}
}
