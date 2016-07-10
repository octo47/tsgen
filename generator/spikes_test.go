package generator

import (
	"fmt"
	"math/rand"

	. "gopkg.in/check.v1"
)

type SpikesTimeSeriesSuite struct {
}

var _ = Suite(&SpikesTimeSeriesSuite{})

func (s *SpikesTimeSeriesSuite) TestSpikes(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewSpikesTimeSeries(rnd)
	points := make([]Point, 100)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		fmt.Println(points[i])
	}
}
