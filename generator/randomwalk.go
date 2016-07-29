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
			rnd:   r,
			lower: lowerBound, upper: upperBound,
			step: step,
			last: Point{Value: r.Float64() / 2 * ((upperBound - lowerBound) + lowerBound)},
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

func (rw *RandomWalkGenerator) UpperBound() float64 {
	return rw.upper
}

func (rw *RandomWalkGenerator) LowerBound() float64 {
	return rw.lower
}

func (rw *RandomWalkGenerator) SetMiddle() {
	rw.last.Value = (rw.upper-rw.lower)/2 + rw.lower
}
