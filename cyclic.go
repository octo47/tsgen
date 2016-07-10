package tsgen

import (
	"math"
	"math/rand"
)

// CyclicTimeSeries implements random walk.
type CyclicTimeSeries struct {
	Generator
	last    Point
	degrees float64
}

func NewCyclicTimeSeries(r *rand.Rand, tags []Tag) *CyclicTimeSeries {
	return &CyclicTimeSeries{Generator{rnd: r, tags: tags}, Point{}, 0}
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
