package generator

import (
	"math"
	"math/rand"

	. "gopkg.in/check.v1"
)

type RandomWalkGeneratorSuite struct {
}

var _ = Suite(&RandomWalkGeneratorSuite{})

func (s *RandomWalkGeneratorSuite) TestRandomWalk(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewRandomWalkGenerator(rnd, 0.01, 0.0, 10.0)
	points := make([]Point, 10)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		diff := points[i+1].Value - points[i].Value
		c.Assert(math.Abs(diff) < 0.5, Equals, true,
			Commentf("Diff should be in range [-0.5;0.5]"))
	}
}
