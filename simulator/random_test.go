package simulator

import (
	"math/rand"
	"sort"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SimulatorSuite struct{}

func TestSimulatorTick(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	conf := NewConfiguration(3, 5, 1)
	simulator := NewSimulator(rnd, conf, 0)
	assert.Equal(t, len(simulator.machines), 3)
	glog.Info(simulator.machines)
	simulator.Run(0, 1, 1600, func(tp *[]TaggedPoints) {
		// nothing to do
	})
}

func TestDeduplicate(t *testing.T) {
	tags := Tags{
		Tag{"tag1", "value1"}, Tag{"tag2", "value2"},
		Tag{"tag2", "value1"}, Tag{"tag1", "value1"},
	}
	sort.Sort(&tags)
	deduplicateTags(&tags)
	dupes := make(map[Tag]bool)
	for ti := range tags {
		_, ok := dupes[tags[ti]]
		assert.False(t, ok)
		dupes[tags[ti]] = true
	}
}

func BenchmarkSimulator(b *testing.B) {
	rnd := rand.New(rand.NewSource(1))
	conf := NewConfiguration(b.N, 1000, 1)
	simulator := NewSimulator(rnd, conf, 0)
	simulator.Run(0, 1, 1600, func(tp *[]TaggedPoints) {
		// nothing to do
	})
}
