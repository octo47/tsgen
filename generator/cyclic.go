package generator

import (
	"math"
	"math/rand"
)

// CyclicGenerator implements random walk.
type CyclicGenerator struct {
	rnd     *rand.Rand
	last    Point
	degrees float64
}

func NewCyclicGenerator(r *rand.Rand) *CyclicGenerator {
	return &CyclicGenerator{rnd: r}
}

func (rw *CyclicGenerator) Next(points *[]Point) {
	for i := range *points {
		rval := rw.rnd.Float64()
		rw.last.Value += rval * math.Sin(rw.degrees*math.Pi/288) / 20
		rw.degrees++
		if rw.degrees > 360 {
			rw.degrees = 0
		}
		(*points)[i] = rw.last
	}
}
