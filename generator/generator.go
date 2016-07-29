package generator

import "math/rand"

type Generator interface {
	// fill slice with generated points
	Next(points *[]Point)
}

type boundedGenerator struct {
	rnd   *rand.Rand
	last  Point
	lower float64
	upper float64
	step  float64
}

type Point struct {
	Timestamp uint64
	Value     float64
}

// calculate bounded step. dv is epxecting to be [0;1)
// step will be scale to boundedGenerator.scale
func calculateNext(dv float64, gen *boundedGenerator) {
	step := gen.step * calculateStep(dv)
	newValue := gen.last.Value + step
	switch {
	case newValue > gen.upper:
		newValue = gen.upper
	case newValue < gen.lower:
		newValue = gen.lower
	}
	gen.last.Value = newValue
}

func calculateStep(dv float64) float64 {
	coeff := 1.0
	if dv < 0.5 {
		coeff = -1.0
	}
	return (dv - 0.5) * (dv - 0.5) * coeff
}
