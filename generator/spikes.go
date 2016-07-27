package generator

import (
	"math"
	"math/rand"
)

// SpikesGenerator implements random walk.
type SpikesGenerator struct {
	rnd      *rand.Rand
	last     Point
	maxSpike float64
}

func NewSpikesGenerator(r *rand.Rand, maxSpike float64) *SpikesGenerator {
	return &SpikesGenerator{rnd: r, maxSpike: maxSpike}
}

func (rw *SpikesGenerator) Next(points *[]Point) {
	for i := range *points {
		spike := math.Ceil(math.Tan(rw.rnd.Float64()*math.Pi/2)) / 100
		if spike > rw.maxSpike {
			spike = rw.maxSpike
		}
		rw.last.Value = spike
		(*points)[i] = rw.last
	}
}
