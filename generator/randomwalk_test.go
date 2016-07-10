package generator

import (
	"math"
	"math/rand"

	. "gopkg.in/check.v1"
)

type RandomWalkTimeSeriesSuite struct {
}

var _ = Suite(&RandomWalkTimeSeriesSuite{})

func (s *RandomWalkTimeSeriesSuite) TestRandomWalk(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewRandomWalkTimeSeries(rnd)
	points := make([]Point, 10)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		diff := points[i+1].Value - points[i].Value
		c.Assert(math.Abs(diff) < 0.5, Equals, true,
			Commentf("Diff should be in range [-0.5;0.5]"))
	}
}
