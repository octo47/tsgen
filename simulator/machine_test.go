package simulator

import (
	"math/rand"

	"github.com/octo47/tsgen/generator"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MachineSuite struct {
	tags []Tag
	rnd  *rand.Rand
}

func newMachineSuite() *MachineSuite {
	return &MachineSuite{
		rnd:  rand.New(rand.NewSource(1)),
		tags: []Tag{Tag{"tag1", "value1"}, Tag{"tag2", "value2"}}}
}

func TestMachineTick(t *testing.T) {
	s := newMachineSuite()
	machine := NewMachine("testhost", s.tags, 0, 2)
	machine.AddTimeseries("service", "requests", generator.NewSpikesGenerator(s.rnd, 100.0), 15)
	cpuTags := []string{"usr", "sys", "io"}
	for _, cpuTag := range cpuTags {
		machine.AddTimeseriesWithTags("sys", "cpu", []Tag{Tag{"cpu.type", cpuTag}},
			generator.NewRandomWalkGenerator(s.rnd, 0.1, 0.0, 100.0), 15)
	}
	machine.AddTimeseries("sys", "disk.usage",
		generator.NewIncreasingGenerator(s.rnd, 0.8, 0.01, 0.1, 0.0, 100.0), 60)
	for timestamp := uint64(300); timestamp < 1200; timestamp += 300 {
		result := machine.Tick(timestamp)
		assert.Equal(t, len(*result), 5)
		for _, taggedPoints := range *result {
			switch *taggedPoints.MetricName {
			case "requests":
				assert.Equal(t, len(*taggedPoints.Datapoints), 300/15)
				assert.Equal(t, int((*taggedPoints.Datapoints)[1].Timestamp-
					(*taggedPoints.Datapoints)[0].Timestamp), 15)
			case "disk.usage":
				assert.Equal(t, len(*taggedPoints.Datapoints), 300/60)
				assert.Equal(t, int((*taggedPoints.Datapoints)[1].Timestamp-
					(*taggedPoints.Datapoints)[0].Timestamp), 60)
			}
		}
	}
}
