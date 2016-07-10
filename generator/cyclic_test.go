package generator

import (
	"math/rand"

	. "gopkg.in/check.v1"
)

type CyclicTimeSeriesSuite struct {
}

var _ = Suite(&CyclicTimeSeriesSuite{})

func (s *CyclicTimeSeriesSuite) TestCyclic(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewCyclicTimeSeries(rnd)
	points := make([]Point, 361)
	rw.Next(&points)
}
