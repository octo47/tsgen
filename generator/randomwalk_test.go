package generator

import (
	"math"
	"math/rand"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomWalk(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewRandomWalkGenerator(rnd, .01, 0.0, 10.0)
	points := make([]Point, 10)
	rw.Next(&points)
	for i := 0; i < len(points)-2; i++ {
		diff := points[i+1].Value - points[i].Value
		assert.True(t, math.Abs(diff) < 0.5,
			"Diff should be in range [-0.5;0.5]")
	}
}
