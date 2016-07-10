package tsgen

import (
	"math/rand"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func TestCyclic(t *testing.T) { TestingT(t) }

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
