package tsgen

import "math/rand"

type TimeSeries interface {
	Tags() []Tag
	// fill slice with generated points
	Next(points *[]Point)
}

type Generator struct {
	tags []Tag
	rnd  *rand.Rand
}

type Tag struct {
	Name  string
	Value string
}

type Point struct {
	Timestamp uint64
	Value     float64
}

func (rw *Generator) Tags() []Tag {
	return rw.tags
}
