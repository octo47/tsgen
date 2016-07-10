package tsgen

import (
	"math"
	"math/rand"

	. "gopkg.in/check.v1"
)

type RandomWalkTimeSeriesSuite struct {
	tags []Tag
}

var _ = Suite(&RandomWalkTimeSeriesSuite{
	tags: []Tag{Tag{"tag1", "value1"}},
})

func (s *RandomWalkTimeSeriesSuite) TestRandomWalk(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewRandomWalkTimeSeries(rnd, s.tags)
	c.Assert(rw.Tags(), DeepEquals, s.tags)
	points := make([]Point, 10)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		diff := points[i+1].Value - points[i].Value
		c.Assert(math.Abs(diff) < 0.5, Equals, true,
			Commentf("Diff should be in range [-0.5;0.5]"))
	}
}
