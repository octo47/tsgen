package tsgen

type TimeSeries interface {
	Tags() []Tag
	// fill slice with generated points
	Next(points *[]Point)
}

type Tag struct {
	Name  string
	Value string
}

type Point struct {
	Timestamp uint64
	Value     float64
}
