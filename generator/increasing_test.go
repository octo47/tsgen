package generator

import (
	"math/rand"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncreasing(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewIncreasingGenerator(rnd, 0.8, 0.01, 0.1, 0.0, 10.0)
	points := make([]Point, 1000)
	rw.Next(&points)
	diffSum := 0.0
	for i := 0; i < len(points)-2; i++ {
		diffSum += points[i+1].Value - points[i].Value
	}
	assert.True(t, diffSum >= 0,
		"Generally we need to be increaseing or at maximum: %v", diffSum)
}
