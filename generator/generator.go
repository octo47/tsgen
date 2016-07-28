package generator

import (
	"math/rand"
)

type Generator interface {
	// fill slice with generated points
	Next(points *[]Point)
}

type boundedGenerator struct {
	rnd   *rand.Rand
	last  Point
	lower float64
	upper float64
}

type Point struct {
	Timestamp uint64
	Value     float64
}

func calculateStep(dv float64) float64 {
	coeff := 1.0
	if dv < 0.5 {
		coeff = -1.0
	}
	return (dv - 0.5) * (dv - 0.5) * coeff
}

func calculateNext(currValue, dv, lower, upper float64) float64 {
	step := calculateStep(dv)
	newValue := currValue + step
	switch {
	case newValue > upper:
		return upper
	case newValue < lower:
		return lower
	}
	return currValue
}
