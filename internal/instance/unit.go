package instance

import (
	"math"
	"strings"
)

type Unit string

const (
	UnitGB = "GB"
	UnitTB = "TB"
)

func (u Unit) BytesToUnit(value float64) float32 {
	switch strings.ToUpper(string(u)) {
	case UnitGB:
		return float32(value / math.Pow(1024, 3))
	case UnitTB:
		return float32(value / math.Pow(1024, 4))
	default:
		return 0
	}
}
