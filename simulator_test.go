package tsgen

import (
	"math/rand"
	"sort"

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
		// nothing to do
	})
}

func (s *SimulatorSuite) TestDeduplicate(c *C) {
	tags := Tags{
		Tag{"tag1", "value1"}, Tag{"tag2", "value2"},
		Tag{"tag2", "value1"}, Tag{"tag1", "value1"},
	}
	sort.Sort(&tags)
	deduplicateTags(&tags)
	dupes := make(map[Tag]bool)
	for ti := range tags {
		_, ok := dupes[tags[ti]]
		c.Assert(ok, Equals, false)
		dupes[tags[ti]] = true
	}
}
