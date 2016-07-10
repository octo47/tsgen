package generator

type TimeSeries interface {
	// fill slice with generated points
	Next(points *[]Point)
}

type Point struct {
	Timestamp uint64
	Value     float64
}
