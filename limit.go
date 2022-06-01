package main

type Unit string

const (
	UnitGB = "GB"
	UnitTB = "TB"
)

type Limit struct {
	Unit
	Value int
}
