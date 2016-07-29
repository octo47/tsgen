package generator

import "math/rand"

// IncreasingGenerator implements random walk.
type IncreasingGenerator struct {
	boundedGenerator
	bias      float64
	resetProb float64
}

func NewIncreasingGenerator(
	r *rand.Rand, bias, resetProb, step, lowerBound, upperBound float64) *IncreasingGenerator {

	return &IncreasingGenerator{
		boundedGenerator: boundedGenerator{
			rnd: r, step: step, lower: lowerBound, upper: upperBound,
		},
		bias:      bias,
		resetProb: resetProb,
	}
}

func (rw *IncreasingGenerator) Next(points *[]Point) {
	for i := range *points {
		if rw.rnd.Float64() < rw.resetProb {
			rw.last.Value = rw.lower
		} else {
			calculateNext(rw.rnd.Float64()+rw.bias, &rw.boundedGenerator)
		}
		(*points)[i] = rw.last
	}
}
