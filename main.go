package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func init() {
	_, err := exec.LookPath("aws")
	if err != nil {
		panic("please install and login aws-cli first")
	}
}

func main() {
	// parse flags
	cPath := flag.String("c", "", "path to json config")
	loop := flag.Int("l", 0, "interval for loop")
	flag.Parse()
	if *cPath == "" {
		panic("config not found")
	}
	// check path
	if _, err := os.Stat(*cPath); err != nil {
		// write empty configuration to path if not exist
		log.Printf("Writing empty config file to %s", *cPath)
		b, _ := json.Marshal(Instances{{
			Name:    "",
			Limit:   Limit{},
			Command: []string{},
		}})
		err = os.WriteFile(*cPath, b, os.ModePerm)
		checkError(err)
		fmt.Println("Config written")
		return
	}
	// load config
	var instances Instances
	b, err := os.ReadFile(*cPath)
	checkError(err)
	err = json.Unmarshal(b, &instances)
	checkError(err)
	instances.Check()
	if *loop > 0 {
		for {
			log.Println("Looping..")
			time.Sleep(time.Duration(*loop) * time.Second)
			instances.Check()
		}
	}
}
