package tsgen

import (
	"github.com/golang/glog"
	"github.com/octo47/tsgen/generator"
)

// Tag is piece of metadata attached to every Point
type Tag struct {
	Name  string
	Value string
}

// TimeSeries keeps track on timestamps and operate on generator to create new points.
type TimeSeries struct {
	name   string
	tags   []Tag
	lastTs uint64
	period uint64
	gen    generator.Generator
}

// TaggedPoints structure hold points with attached Tags
type TaggedPoints struct {
	name   *string
	tags   *[]Tag
	points *[]generator.Point
}

// Machine sturct keeps state of simulated machine, including timeseries and metadata
type Machine struct {
	name       string
	tags       []Tag
	timeseries []TimeSeries
	lastTs     uint64
}

// NewMachine creates Machine with specified tags and initial timestamp.
func NewMachine(name string, tags []Tag, initialTs uint64) *Machine {
	return &Machine{name: name, tags: tags, lastTs: initialTs}
}

// AddTimeseries adds timeseries with generator and period.
func (m *Machine) AddTimeseries(name string, gen generator.Generator, period uint64) {
	m.timeseries = append(m.timeseries,
		TimeSeries{name: name, tags: m.tags, lastTs: m.lastTs, period: period, gen: gen})
}

// AddTimeseriesWithTags adds timeseries with series specific tags with generator and period.
func (m *Machine) AddTimeseriesWithTags(
	name string, tags []Tag, gen generator.Generator, period uint64) {
	m.timeseries = append(m.timeseries,
		TimeSeries{name: name, tags: append(tags, m.tags...), lastTs: m.lastTs, period: period, gen: gen})
}

// Tick advance time for machine. Machine generates all metrics up to provided timestamp.
// Returns tagged points for every timeseries.
func (m *Machine) Tick(timestamp uint64) *[]TaggedPoints {
	tpoints := make([]TaggedPoints, len(m.timeseries))
	for tsidx := range m.timeseries {
		tpoints[tsidx] = m.timeseries[tsidx].Tick(timestamp)
	}
	return &tpoints
}

// Tick advaces single timeseries to timestamp, generating points in between last timestamp and new
func (ts *TimeSeries) Tick(timestamp uint64) TaggedPoints {
	count := (timestamp - ts.lastTs) / ts.period
	if count < 1 {
		return TaggedPoints{&ts.name, &ts.tags, nil}
	}
	buffer := make([]generator.Point, count)
	ts.gen.Next(&buffer)
	lastTimestamp := ts.lastTs
	for bidx := range buffer {
		lastTimestamp += ts.period
		buffer[bidx].Timestamp = lastTimestamp
	}
	ts.lastTs = lastTimestamp
	if glog.V(2) {
		glog.Info("generated ", ts.name, " tick, ts=", timestamp, " points=", len(buffer),
			" timesries=", ts)
	}
	return TaggedPoints{&ts.name, &ts.tags, &buffer}
}
