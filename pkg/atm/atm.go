package atm

import (
	"context"
	"fmt"
	"github.com/SuzukiHonoka/AWS-Traffic-Monitor/internal/api"
	"log"
	"time"
)

type ATM struct {
	cfg    *Config
	cancel func()
}

func New(cfg *Config) (*ATM, error) {
	// Sanity check
	if cfg.AWSConfig == nil {
		return nil, ErrAWSConfigNotConfigured
	}

	// Init api
	if err := api.Init(cfg.AWSConfig); err != nil {
		return nil, fmt.Errorf("failed to init api, err=%v", err)
	}

	atm := &ATM{cfg: cfg}
	return atm, nil
}

func (a *ATM) Start() error {
	a.cfg.InstanceList.Check()

	if a.cfg.LoopIntervalString == "" {
		return nil
	}
	d, err := time.ParseDuration(a.cfg.LoopIntervalString)
	if err != nil {
		return fmt.Errorf("failed to parse loop interval, err=%v", err)
	}

	// looping
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			log.Println("Looping..")
			time.Sleep(d)
			a.cfg.InstanceList.Check()
		}
	}
}

func (a *ATM) Stop() {
	if a.cancel != nil {
		log.Println("Stopping ATM...")
		a.cancel()
	}
}
