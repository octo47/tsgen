package generator

import (
	"math/rand"

	"testing"
)

func TestSpikes(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	rw := NewSpikesGenerator(rnd, 5.0)
	points := make([]Point, 20)
	rw.Next(&points)
}
