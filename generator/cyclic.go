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
		boundedGenerator: boundedGenerator{
			rnd: r, step: (upperBound - lowerBound) / 100.0,
			lower: lowerBound, upper: upperBound,
			last: Point{Value: r.Float64() * ((upperBound - lowerBound) + lowerBound)},
		}}
}

func (rw *CyclicGenerator) Next(points *[]Point) {
	for i := range *points {
		calculateNext(
			(rw.rnd.Float64() + math.Sin(float64(rw.counter)*math.Pi/288)/100),
			&rw.boundedGenerator)
		rw.counter++
		(*points)[i] = rw.last
	}
}

func (rw *CyclicGenerator) UpperBound() float64 {
	return rw.upper
}

func (rw *CyclicGenerator) LowerBound() float64 {
	return rw.lower
}

func (rw *CyclicGenerator) SetMiddle() {
	rw.last.Value = (rw.upper-rw.lower)/2 + rw.lower
}
