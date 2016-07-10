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
		rvalue := rw.rnd.Float64()
		coeff := float64(1)
		if rvalue < 0.5 {
			coeff = -1
		}
		rw.last.Value += (rvalue - 0.5) * (rvalue - 0.5) * coeff
		(*points)[i] = rw.last
	}
}
