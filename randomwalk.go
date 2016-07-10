package tsgen

import "math/rand"

// RandomWalkTimeSeries implements random walk.
type RandomWalkTimeSeries struct {
	last Point
	tags []Tag
	rnd  *rand.Rand
}

func NewRandomWalkTimeSeries(r *rand.Rand, tags []Tag) *RandomWalkTimeSeries {
	return &RandomWalkTimeSeries{rnd: r, tags: tags}
}

func (rw *RandomWalkTimeSeries) Tags() []Tag {
	return rw.tags
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
