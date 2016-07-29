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
		spike *= rw.upper
		if rw.last.Value > spike {
			rw.last.Value *= rw.rnd.Float64() // decay spike
		} else {
			rw.last.Value = spike
		}
		(*points)[i] = rw.last
	}
}
