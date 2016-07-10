package tsgen

import (
	"fmt"
	"math/rand"

	. "gopkg.in/check.v1"
)

type SpikesTimeSeriesSuite struct {
	tags []Tag
}

var _ = Suite(&SpikesTimeSeriesSuite{
	tags: []Tag{Tag{"tag1", "value1"}},
})

func (s *SpikesTimeSeriesSuite) TestSpikes(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewSpikesTimeSeries(rnd, s.tags)
	c.Assert(rw.Tags(), DeepEquals, s.tags)
	points := make([]Point, 100)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		fmt.Println(points[i])
	}
}
