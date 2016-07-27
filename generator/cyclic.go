package generator

import (
	"math"
	"math/rand"
)

// CyclicGenerator implements random walk.
type CyclicGenerator struct {
	rnd     *rand.Rand
	last    Point
	bias    float64
	counter int
}

func NewCyclicGenerator(r *rand.Rand) *CyclicGenerator {
	return &CyclicGenerator{rnd: r, counter: r.Intn(365)}
}

func (rw *CyclicGenerator) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value = calculateNext(
			rw.last.Value,
			rw.rnd.Float64()+math.Sin(float64(rw.counter)*math.Pi/288)/20)
		rw.counter++
		(*points)[i] = rw.last
	}
}
