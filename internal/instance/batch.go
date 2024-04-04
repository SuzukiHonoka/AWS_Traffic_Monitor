package instance

import (
	"encoding/json"
	"os"
)

type Instances []Instance

func Load(path string) (*Instances, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var instances Instances
	if err = json.Unmarshal(b, &instances); err != nil {
		return nil, err
	}
	return &instances, nil
}

func (s Instances) Check() {
	for _, v := range s {
		v.Check()
	}
}
