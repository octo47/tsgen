package tsgen

import (
	"math"
	"math/rand"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type RandomWalkTimeSeriesSuite struct {
	tags []Tag
}

var _ = Suite(&RandomWalkTimeSeriesSuite{
	tags: []Tag{Tag{"tag1", "value1"}},
})

func (s *RandomWalkTimeSeriesSuite) TestHelloWorld(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewRandomWalkTimeSeries(rnd, s.tags)
	c.Assert(rw.Tags(), DeepEquals, s.tags)
	points := make([]Point, 10)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		diff := points[i].Value - points[i+1].Value
		c.Assert(math.Abs(diff) < 0.5, Equals, true)
	}
}
