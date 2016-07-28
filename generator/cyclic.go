package generator

import (
	"math"
	"math/rand"
)

// CyclicGenerator implements random walk.
type CyclicGenerator struct {
	boundedGenerator
	counter int
}

func NewCyclicGenerator(r *rand.Rand, lowerBound, upperBound float64) *CyclicGenerator {
	return &CyclicGenerator{counter: r.Intn(365),
		boundedGenerator: boundedGenerator{rnd: r, lower: lowerBound, upper: upperBound}}
}

func (rw *CyclicGenerator) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value = calculateNext(
			rw.last.Value,
			rw.rnd.Float64()+math.Sin(float64(rw.counter)*math.Pi/288)/20,
			rw.lower,
			rw.upper)
		rw.counter++
		(*points)[i] = rw.last
	}
}
