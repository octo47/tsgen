package generator

import (
	"math/rand"

	. "gopkg.in/check.v1"
)

type IncreasingGeneratorSuite struct {
}

var _ = Suite(&IncreasingGeneratorSuite{})

func (s *IncreasingGeneratorSuite) TestIncreasing(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewIncreasingGenerator(rnd, 0.8, 0.01, 0.1, 0.0, 10.0)
	points := make([]Point, 1000)
	rw.Next(&points)
	diffSum := 0.0
	for i := 0; i < len(points)-2; i++ {
		diffSum += points[i+1].Value - points[i].Value
	}
	c.Assert(diffSum >= 0, Equals, true,
		Commentf("Generally we need to be increaseing or at maximum: %v", diffSum))
}
