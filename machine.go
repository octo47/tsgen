package tsgen

import (
	"strings"

	"github.com/golang/glog"
	"github.com/octo47/tsgen/generator"
)

// Tag is piece of metadata attached to every Point
type Tag struct {
	Name  string
	Value string
}

// Tags is for sorting here.
type Tags []Tag

// TimeSeries keeps track on timestamps and operate on generator to create new points.
type TimeSeries struct {
	ns     string
	name   string
	tags   Tags
	lastTs uint64
	period uint64
	gen    generator.Generator
}

// TaggedPoints structure hold points with attached Tags
type TaggedPoints struct {
	Namespace  *string
	MetricName *string
	Tags       *Tags
	Datapoints *[]generator.Point
}

// Machine sturct keeps state of simulated machine, including timeseries and metadata
type Machine struct {
	name       string
	tags       Tags
	timeseries []TimeSeries
	lastTs     uint64
}

// NewMachine creates Machine with specified tags and initial timestamp.
func NewMachine(name string, tags Tags, initialTs uint64) *Machine {
	return &Machine{name: name, tags: tags, lastTs: initialTs}
}

// AddTimeseries adds timeseries with generator and period.
func (m *Machine) AddTimeseries(ns string, name string, gen generator.Generator, period uint64) {
	m.timeseries = append(m.timeseries,
		TimeSeries{ns: ns, name: name, tags: m.tags, lastTs: m.lastTs, period: period, gen: gen})
}

// AddTimeseriesWithTags adds timeseries with series specific tags with generator and period.
func (m *Machine) AddTimeseriesWithTags(
	ns string, name string, tags Tags, gen generator.Generator, period uint64) {
	m.timeseries = append(m.timeseries,
		TimeSeries{ns: ns, name: name, tags: append(tags, m.tags...), lastTs: m.lastTs, period: period, gen: gen})
}

// Tick advance time for machine. Machine generates all metrics up to provided timestamp.
// Returns tagged points for every timeseries.
func (m *Machine) Tick(timestamp uint64) *[]TaggedPoints {
	if glog.V(2) {
		glog.Info(" machine ", m.name, " ", m.lastTs, " -> ", timestamp)
	}
	tpoints := make([]TaggedPoints, len(m.timeseries))
	for tsidx := range m.timeseries {
		tpoints[tsidx] = m.timeseries[tsidx].Tick(timestamp)
	}
	return &tpoints
}

// Tick advaces single timeseries to timestamp, generating points in between last timestamp and new
func (ts *TimeSeries) Tick(timestamp uint64) TaggedPoints {
	count := (timestamp - ts.lastTs) / ts.period
	if glog.V(2) {
		glog.Info(" timeseries: ", ts.name, " ", ts.lastTs, " -> ", timestamp, " ", count)
	}
	if count < 1 {
		return TaggedPoints{&ts.ns, &ts.name, &ts.tags, nil}
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
	return TaggedPoints{&ts.ns, &ts.name, &ts.tags, &buffer}
}

func (tags *Tags) Len() int {
	return len(*tags)
}

func (tags *Tags) Less(i int, j int) bool {
	c := strings.Compare((*tags)[i].Name, (*tags)[j].Name)
	if c == -1 {
		return true
	}
	c = strings.Compare((*tags)[i].Value, (*tags)[j].Value)
	return c == -1
}

func (tags *Tags) Swap(i int, j int) {
	(*tags)[i], (*tags)[j] = (*tags)[j], (*tags)[i]
}
