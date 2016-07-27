package generator

import (
	"math"
	"math/rand"
)

// SpikesGenerator implements random walk.
type SpikesGenerator struct {
	rnd  *rand.Rand
	last Point
}

func NewSpikesGenerator(r *rand.Rand) *SpikesGenerator {
	return &SpikesGenerator{rnd: r}
}

func (rw *SpikesGenerator) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value = math.Ceil(math.Tan(rw.rnd.Float64()*math.Pi/2)) / 1000
		(*points)[i] = rw.last
	}
}
