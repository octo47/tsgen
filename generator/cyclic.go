package generator

import (
	"math"
	"math/rand"
)

// CyclicTimeSeries implements random walk.
type CyclicTimeSeries struct {
	rnd     *rand.Rand
	last    Point
	degrees float64
}

func NewCyclicTimeSeries(r *rand.Rand) *CyclicTimeSeries {
	return &CyclicTimeSeries{rnd: r}
}

func (rw *CyclicTimeSeries) Next(points *[]Point) {
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
