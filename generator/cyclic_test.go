package generator

import (
	"fmt"
	"math/rand"

	. "gopkg.in/check.v1"
)

type CyclicGeneratorSuite struct {
}

var _ = Suite(&CyclicGeneratorSuite{})

func (s *CyclicGeneratorSuite) TestCyclic(c *C) {
	rnd := rand.New(rand.NewSource(1234))
	rw := NewCyclicGenerator(rnd, 0.0, 100.0)
	points := make([]Point, 512)
	rw.Next(&points)
	fmt.Println(points)
}
