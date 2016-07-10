package tsgen

import "math/rand"

// IncreasingTimeSeries implements random walk.
type IncreasingTimeSeries struct {
	Generator
	last Point
}

func NewIncreasingTimeSeries(r *rand.Rand, tags []Tag) *IncreasingTimeSeries {
	return &IncreasingTimeSeries{Generator{rnd: r, tags: tags}, Point{}}
}

func (rw *IncreasingTimeSeries) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value += rw.rnd.Float64()
		(*points)[i] = rw.last
	}
}
