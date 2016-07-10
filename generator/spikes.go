package generator

import (
	"math"
	"math/rand"
)

// SpikesTimeSeries implements random walk.
type SpikesTimeSeries struct {
	rnd  *rand.Rand
	last Point
}

func NewSpikesTimeSeries(r *rand.Rand) *SpikesTimeSeries {
	return &SpikesTimeSeries{rnd: r}
}

func (rw *SpikesTimeSeries) Next(points *[]Point) {
	for i := range *points {
		rw.last.Value = math.Ceil(math.Tan(rw.rnd.Float64()*math.Pi/2)) / 100
		(*points)[i] = rw.last
	}
}
