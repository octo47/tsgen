package tsgen

import (
	"math/rand"

	. "gopkg.in/check.v1"
)

type CyclicTimeSeriesSuite struct {
	tags []Tag
}

var _ = Suite(&CyclicTimeSeriesSuite{
	tags: []Tag{Tag{"tag1", "value1"}},
})

func (s *CyclicTimeSeriesSuite) TestCyclic(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewCyclicTimeSeries(rnd, s.tags)
	c.Assert(rw.Tags(), DeepEquals, s.tags)
	points := make([]Point, 361)
	rw.Next(&points)
}
