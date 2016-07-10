package tsgen

import (
	"math/rand"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func TestIncreasing(t *testing.T) { TestingT(t) }

type IncreasingTimeSeriesSuite struct {
	tags []Tag
}

var _ = Suite(&IncreasingTimeSeriesSuite{
	tags: []Tag{Tag{"tag1", "value1"}},
})

func (s *IncreasingTimeSeriesSuite) TestIncreasing(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewIncreasingTimeSeries(rnd, s.tags)
	c.Assert(rw.Tags(), DeepEquals, s.tags)
	points := make([]Point, 10)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		diff := points[i+1].Value - points[i].Value
		c.Assert(diff >= 0, Equals, true, Commentf("Diff should be positive: %v", diff))
	}
}
