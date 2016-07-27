package generator

import "math/rand"

// RandomWalkGenerator implements random walk.
type RandomWalkGenerator struct {
	rnd  *rand.Rand
	last Point
}

func NewRandomWalkGenerator(r *rand.Rand) *RandomWalkGenerator {
	return &RandomWalkGenerator{rnd: r}
}

func (rw *RandomWalkGenerator) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value = calculateNext(rw.last.Value, rw.rnd.Float64())
		(*points)[i] = rw.last
	}
}
