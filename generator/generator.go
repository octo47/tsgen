package generator

import "math"

type Generator interface {
	// fill slice with generated points
	Next(points *[]Point)
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

func calculateNext(currValue float64, dv float64) float64 {
	step := calculateStep(dv)
	// make ajustements only if we are going to 'reduce' postivie value
	if step < 0.0 && currValue+step < 0.0 {
		// taking absolute value and using it for calculating relative decrease
		// but keeping value positive
		currValue *= math.Abs(currValue) / math.Abs(step)
	} else {
		currValue += step
	}
	return currValue
}
