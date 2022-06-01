package main

import (
	"errors"
	"math"
	"os/exec"
	"strings"
	"time"
)

func BeginningOfDay(now time.Time) time.Time {
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

func EndOfDay(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
}

func BytesToUnit(unit Unit, value float64) float32 {
	switch strings.ToUpper(string(unit)) {
	case UnitGB:
		return float32(value / math.Pow(1024, 3))
	case UnitTB:
		return float32(value / math.Pow(1024, 4))
	default:
		panic(errors.New("unit not supported"))
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Exec(cmd string) []byte {
	args := strings.Split(cmd, " ")
	result, err := exec.Command(args[0], args[1:]...).Output()
	checkError(err)
	return result
}
