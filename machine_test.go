package tsgen

import (
	"math/rand"

	"github.com/octo47/tsgen/generator"
	. "gopkg.in/check.v1"
)

type MachineSuite struct {
	tags []Tag
	rnd  *rand.Rand
}

var _ = Suite(&MachineSuite{
	rnd:  rand.New(rand.NewSource(1)),
	tags: []Tag{Tag{"tag1", "value1"}, Tag{"tag2", "value2"}},
})

func (s *MachineSuite) TestMachineTick(c *C) {
	machine := NewMachine("testhost", s.tags, 0)
	machine.AddTimeseries("requests", generator.NewSpikesGenerator(s.rnd), 15)
	cpuTags := []string{"usr", "sys", "io"}
	for _, cpuTag := range cpuTags {
		machine.AddTimeseriesWithTags("cpu", []Tag{Tag{"cpu.type", cpuTag}},
			generator.NewRandomWalkGenerator(s.rnd), 15)
	}
	machine.AddTimeseries("disk.usage", generator.NewIncreasingGenerator(s.rnd), 60)
	for timestamp := uint64(300); timestamp < 1200; timestamp += 300 {
		result := machine.Tick(timestamp)
		c.Assert(len(*result), Equals, 5)
		for _, taggedPoints := range *result {
			switch *taggedPoints.name {
			case "requests":
				c.Assert(len(*taggedPoints.points), Equals, 300/15)
				c.Assert(int((*taggedPoints.points)[1].Timestamp-(*taggedPoints.points)[0].Timestamp),
					Equals, 15)
			case "disk.usage":
				c.Assert(len(*taggedPoints.points), Equals, 300/60)
				c.Assert(int((*taggedPoints.points)[1].Timestamp-(*taggedPoints.points)[0].Timestamp),
					Equals, 60)
			}
		}
	}
}
