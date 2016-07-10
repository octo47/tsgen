package tsgen

import "math/rand"

// RandomWalkTimeSeries implements random walk.
type RandomWalkTimeSeries struct {
	Generator
	last Point
}

func NewRandomWalkTimeSeries(r *rand.Rand, tags []Tag) *RandomWalkTimeSeries {
	return &RandomWalkTimeSeries{Generator{rnd: r, tags: tags}, Point{}}
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
