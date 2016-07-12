package generator

import (
	"math"
	"math/rand"
)

// RandomWalkGenerator implements random walk.
type RandomWalkGenerator struct {
	rnd      *rand.Rand
	last     Point
	positive bool
}

func NewRandomWalkGenerator(r *rand.Rand) *RandomWalkGenerator {
	return &RandomWalkGenerator{rnd: r}
}

func NewRandomWalkGeneratorPositive(r *rand.Rand) *RandomWalkGenerator {
	return &RandomWalkGenerator{rnd: r, positive: true}
}

func (rw *RandomWalkGenerator) Next(points *[]Point) {
	for i := range *points {
		rvalue := rw.rnd.Float64()
		coeff := float64(1)
		if rvalue < 0.5 {
			coeff = -1
		}
		step := (rvalue - 0.5) * (rvalue - 0.5) * coeff
		rw.last.Value += (rvalue - 0.5) * (rvalue - 0.5) * coeff
		if rw.positive && rw.last.Value < 0.0 {
			rw.last.Value += 2 * math.Abs(step)
		}
		(*points)[i] = rw.last
	}
}
