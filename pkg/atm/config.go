package atm

import (
	"AWS_Trafiic_Monitor/internal/api"
	"AWS_Trafiic_Monitor/internal/instance"
	"encoding/json"
	"os"
)

type Config struct {
	LoopIntervalString string        `json:"interval"`
	InstanceList       instance.List `json:"instance"`
	AWSConfig          *api.Config   `json:"aws"`
}

func LoadFromFile(path string) (*Config, error) {
	// Read instance config file
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into Config
	c := new(Config)
	if err = json.Unmarshal(b, c); err != nil {
		return nil, err
	}

	return c, nil
}
