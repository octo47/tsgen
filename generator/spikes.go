package generator

import (
	"math"
	"math/rand"
)

// SpikesGenerator implements random walk.
type SpikesGenerator struct {
	rnd   *rand.Rand
	upper float64
	last  Point
}

func NewSpikesGenerator(r *rand.Rand, upperBound float64) *SpikesGenerator {
	return &SpikesGenerator{rnd: r, upper: upperBound}
}

func (rw *SpikesGenerator) Next(points *[]Point) {
	for i := range *points {
		spike := math.Ceil(math.Tan(rw.rnd.Float64()*math.Pi/2)) / 1000.0
		if spike > 1.0 {
			// tan can be quite large
			spike = 1.0
		}
		if spike < 0.5 {
			spike = 0
		}
		spike *= rw.upper
		if rw.last.Value > spike {
			rw.last.Value *= rw.rnd.Float64() // decay spike
		} else {
			rw.last.Value = spike
		}
		(*points)[i] = rw.last
	}
}

func (rw *SpikesGenerator) UpperBound() float64 {
	return rw.upper
}

func (rw *SpikesGenerator) LowerBound() float64 {
	return 0.0
}

func (rw *SpikesGenerator) SetMiddle() {
	// we don't support middle for spikes, so just reset to zero
	rw.last.Value = 0.0
}
