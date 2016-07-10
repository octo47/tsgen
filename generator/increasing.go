package generator

import "math/rand"

// IncreasingTimeSeries implements random walk.
type IncreasingTimeSeries struct {
	rnd  *rand.Rand
	last Point
}

func NewIncreasingTimeSeries(r *rand.Rand) *IncreasingTimeSeries {
	return &IncreasingTimeSeries{rnd: r}
}

func (rw *IncreasingTimeSeries) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value += rw.rnd.Float64()
		(*points)[i] = rw.last
	}
}
