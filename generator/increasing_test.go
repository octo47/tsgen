package generator

import (
	"math/rand"

	. "gopkg.in/check.v1"
)

type IncreasingTimeSeriesSuite struct {
}

var _ = Suite(&IncreasingTimeSeriesSuite{})

func (s *IncreasingTimeSeriesSuite) TestIncreasing(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewIncreasingTimeSeries(rnd)
	points := make([]Point, 10)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		diff := points[i+1].Value - points[i].Value
		c.Assert(diff >= 0, Equals, true, Commentf("Diff should be positive: %v", diff))
	}
}
