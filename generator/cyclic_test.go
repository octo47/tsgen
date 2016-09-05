package generator

import (
	"math/rand"
	"testing"
)

func TestCyclic(t *testing.T) {
	rnd := rand.New(rand.NewSource(1234))
	rw := NewCyclicGenerator(rnd, 0.0, 100.0)
	points := make([]Point, 512)
	rw.Next(&points)
}
