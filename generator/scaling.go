package generator

// ScalingGenerator implements random walk.
type ScalingGenerator struct {
	inner Generator
	scale float64
}

func NewScalingGenerator(inner Generator, scale float64) *ScalingGenerator {
	return &ScalingGenerator{inner: inner, scale: scale}
}

func (rw *ScalingGenerator) Next(points *[]Point) {
	rw.inner.Next(points)
	for i := range *points {
		(*points)[i].Value *= rw.scale
	}
}

func (rw *ScalingGenerator) UpperBound() float64 {
	return rw.inner.UpperBound() * rw.scale
}

func (rw *ScalingGenerator) LowerBound() float64 {
	return rw.inner.LowerBound() * rw.scale
}

func (rw *ScalingGenerator) SetMiddle() {
	rw.inner.SetMiddle()
}
