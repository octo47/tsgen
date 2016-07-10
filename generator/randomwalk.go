package generator

import "math/rand"

// RandomWalkTimeSeries implements random walk.
type RandomWalkTimeSeries struct {
	rnd  *rand.Rand
	last Point
}

func NewRandomWalkTimeSeries(r *rand.Rand) *RandomWalkTimeSeries {
	return &RandomWalkTimeSeries{rnd: r}
}

func (rw *RandomWalkTimeSeries) Next(points *[]Point) {
	for i := range *points {
		rvalue := rw.rnd.Float64()
		coeff := float64(1)
		if rvalue < 0.5 {
			coeff = -1
		}
		rw.last.Value += (rvalue - 0.5) * (rvalue - 0.5) * coeff
		(*points)[i] = rw.last
	}
}
