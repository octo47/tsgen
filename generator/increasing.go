package generator

import "math/rand"

// IncreasingGenerator implements random walk.
type IncreasingGenerator struct {
	boundedGenerator
	bias      float64
	resetProb float64
}

func NewIncreasingGenerator(
	r *rand.Rand, bias, resetProb, lowerBound, upperBound float64) *IncreasingGenerator {

	return &IncreasingGenerator{
		boundedGenerator: boundedGenerator{rnd: r, lower: lowerBound, upper: upperBound},
		bias:             bias,
		resetProb:        resetProb,
	}
}

func (rw *IncreasingGenerator) Next(points *[]Point) {
	for i := range *points {
		if rw.rnd.Float64() < rw.resetProb {
			rw.last.Value = rw.lower
		} else {
			rw.last.Value = calculateNext(rw.last.Value, rw.rnd.Float64()+rw.bias,
				rw.lower, rw.upper)
		}
		(*points)[i] = rw.last
	}
}
