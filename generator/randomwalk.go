package generator

import "math/rand"

// RandomWalkGenerator implements random walk.
type RandomWalkGenerator struct {
	boundedGenerator
	step float64
}

func NewRandomWalkGenerator(
	r *rand.Rand, step, lowerBound, upperBound float64) *RandomWalkGenerator {

	return &RandomWalkGenerator{
		boundedGenerator: boundedGenerator{
			rnd: r, lower: lowerBound, upper: upperBound,
		},
		step: step,
	}
}

func (rw *RandomWalkGenerator) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value = calculateNext(rw.last.Value, rw.step*rw.rnd.Float64(), rw.lower, rw.upper)
		(*points)[i] = rw.last
	}
}
