package generator

// CombineGenerator uses one generator as base and other generators as mutators.
type CombineGenerator struct {
	base     Generator
	mutators []Generator
}

func NewCombineGenerator(base Generator, mutators []Generator) *CombineGenerator {
	return &CombineGenerator{base: base, mutators: mutators}
}

func (rw *CombineGenerator) Next(points *[]Point) {
	rw.base.Next(points)
	addPoints := make([]Point, len(*points))
	for mi := range rw.mutators {
		rw.mutators[mi].Next(&addPoints)
		for i, add := range addPoints {
			newValue := (*points)[i].Value + add.Value
			switch {
			case newValue > rw.base.UpperBound():
				newValue = rw.base.UpperBound()
				rw.mutators[mi].SetMiddle()
			case newValue < rw.base.LowerBound():
				newValue = rw.base.LowerBound()
				rw.mutators[mi].SetMiddle()
			}
			(*points)[i].Value = newValue
		}
	}
}

func (rw *CombineGenerator) UpperBound() float64 {
	return rw.base.UpperBound()
}

func (rw *CombineGenerator) LowerBound() float64 {
	return rw.base.LowerBound()
}

func (rw *CombineGenerator) SetMiddle() {
	rw.base.SetMiddle()
	for mi := range rw.mutators {
		rw.mutators[mi].SetMiddle()
	}
}
