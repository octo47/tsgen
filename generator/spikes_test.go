package generator

import (
	"math/rand"

	. "gopkg.in/check.v1"
)

type SpikesGeneratorSuite struct {
}

var _ = Suite(&SpikesGeneratorSuite{})

func (s *SpikesGeneratorSuite) TestSpikes(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewSpikesGenerator(rnd)
	points := make([]Point, 20)
	rw.Next(&points)
}
