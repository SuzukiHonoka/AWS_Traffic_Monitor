package main

import (
	"flag"
	"github.com/SuzukiHonoka/AWS-Traffic-Monitor/pkg/atm"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	// General config
	configPath = flag.String("c", "./config.json", "path to json config")
)

func init() {
	// parse flags
	flag.Parse()
}

func main() {
	// Load config
	cfg, err := atm.LoadFromFile(*configPath)
	if err != nil {
		log.Fatalf("failed to load config, err=%v", err)
	}

	// Create ATM instance
	a, err := atm.New(cfg)
	if err != nil {
		log.Fatalf("failed to create ATM instance, err=%v", err)
	}

	// Start ATM
	go func() {
		if err = a.Start(); err != nil {
			log.Fatalf("failed to start ATM, err=%v", err)
		}
	}()
	defer a.Stop()

	// Wait for interrupt signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Received interrupt signal, stopping ATM...")
}
