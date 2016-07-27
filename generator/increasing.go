package generator

import "math/rand"

// IncreasingGenerator implements random walk.
type IncreasingGenerator struct {
	rnd  *rand.Rand
	last Point
	bias float64
}

func NewIncreasingGenerator(r *rand.Rand, bias float64) *IncreasingGenerator {
	return &IncreasingGenerator{rnd: r}
}

func (rw *IncreasingGenerator) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value = calculateNext(rw.last.Value, rw.rnd.Float64()+rw.bias)
		(*points)[i] = rw.last
	}
}
