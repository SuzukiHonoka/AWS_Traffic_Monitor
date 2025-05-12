package instance

import (
	"strings"
)

type Unit string

const (
	UnitGB = "GB"
	UnitTB = "TB"
)

func (u Unit) BytesToUnit(value float64) float64 {
	switch strings.ToUpper(string(u)) {
	case UnitGB:
		return value / (1 << 30)
	case UnitTB:
		return value / (1 << 40)
	default:
		return 0
	}
}
