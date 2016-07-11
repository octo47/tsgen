package generator

import "math/rand"

// IncreasingGenerator implements random walk.
type IncreasingGenerator struct {
	rnd  *rand.Rand
	last Point
}

func NewIncreasingGenerator(r *rand.Rand) *IncreasingGenerator {
	return &IncreasingGenerator{rnd: r}
}

func (rw *IncreasingGenerator) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value += rw.rnd.Float64()
		(*points)[i] = rw.last
	}
}