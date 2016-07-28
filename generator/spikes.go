package generator

import (
	"math"
	"math/rand"
)

// SpikesGenerator implements random walk.
type SpikesGenerator struct {
	boundedGenerator
}

func NewSpikesGenerator(r *rand.Rand, upperBound float64) *SpikesGenerator {
	return &SpikesGenerator{boundedGenerator{rnd: r, upper: upperBound, lower: 0.0}}
}

func (rw *SpikesGenerator) Next(points *[]Point) {
	for i := range *points {
		spike := math.Ceil(math.Tan(rw.rnd.Float64()*math.Pi/2)) / 100
		if spike > rw.upper {
			spike = rw.upper
		}
		rw.last.Value = spike
		(*points)[i] = rw.last
	}
}
