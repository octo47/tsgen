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
			rnd: r, lower: lowerBound, upper: upperBound, step: step,
		},
		step: step,
	}
}

func (rw *RandomWalkGenerator) Next(points *[]Point) {
	for i := range *points {
		calculateNext(rw.rnd.Float64(), &rw.boundedGenerator)
		(*points)[i] = rw.last
	}
}
