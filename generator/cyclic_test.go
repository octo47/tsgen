package generator

import (
	"math/rand"

	. "gopkg.in/check.v1"
)

type CyclicGeneratorSuite struct {
}

var _ = Suite(&CyclicGeneratorSuite{})

func (s *CyclicGeneratorSuite) TestCyclic(c *C) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewCyclicGenerator(rnd)
	points := make([]Point, 361)
	rw.Next(&points)
}