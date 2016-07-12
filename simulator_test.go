package tsgen

import (
	"fmt"
	"math/rand"

	"github.com/golang/glog"

	. "gopkg.in/check.v1"
)

type SimulatorSuite struct {
	tags []Tag
	rnd  *rand.Rand
}

var _ = Suite(&SimulatorSuite{
	rnd:  rand.New(rand.NewSource(1)),
	tags: []Tag{Tag{"tag1", "value1"}, Tag{"tag2", "value2"}},
})

func (s *SimulatorSuite) TestSimulatorTick(c *C) {
	rnd := rand.New(rand.NewSource(1))
	conf := NewConfiguration(3, 5)
	simulator := NewSimulator(rnd, conf, 0)
	c.Assert(len(simulator.machines), Equals, 3)
	glog.Info(simulator.machines)
	simulator.Run(0, 1, 1600, func(tp *[]TaggedPoints) {
		fmt.Println(*tp)
	})
}
