package simulator

import (
	"bytes"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/golang/glog"
	"github.com/octo47/tsgen/generator"
)

// Simulator keeps state of simulation.
type Simulator struct {
	conf        Configuration
	currentTime uint64
	machines    []*Machine
	tags        tagsDef
}

// Configuration keeps simulation configuration
type Configuration struct {
	// total machines count
	Machines int
	// group machines by clusters.
	// machine in every cluster will drain unique tags from cluster private set of possible tags
	Clusters int
	// number of tags that will be assigned to every machine
	// value will be unique to machine
	GlobalTags int
	// unique tags are tags assigned per machine
	ClusterTags int
	// minimun tags
	MinimumTags int
	// Total metrics
	MetricsTotal int
	// Set of base metrics that reported by all machines
	BaseMetrics int
	// Maximum metrics per machine
	MaxMetrics int
	// Metrics per namespace
	MetricsPerNamespace int
	// Tags per metric
	TagsPerMetric int
	// TagValuesPerMetricTag configures how deiverse values of metric tags.
	TagValuesPerMetricTag int
	// StartSplay configures how far machines start time could drift away from startTime
	StartSplay int
}

func NewConfiguration(machines int, metrics int, clusters int) Configuration {
	return Configuration{
		Machines:              machines,
		Clusters:              clusters,
		GlobalTags:            8,
		ClusterTags:           clusters/48 + 16, // suppose 48 machines average per service
		MinimumTags:           4,
		MetricsTotal:          int(math.Log10(float64(machines)) * 300),
		BaseMetrics:           metrics/10 + 1,
		MaxMetrics:            metrics,
		MetricsPerNamespace:   metrics/100 + 50,
		TagsPerMetric:         2,
		TagValuesPerMetricTag: 64,
		StartSplay:            300,
	}
}

type clusterDef struct {
	cID     int
	tags    tagsDef
	metrics []metricDef
}

func NewSimulator(rnd *rand.Rand, conf Configuration, startTime uint64) *Simulator {
	clusters := generateClusters(rnd, conf)
	if glog.V(1) {
		glog.Info("Generated ", len(clusters), " clusters")
	}
	clusterGen := rand.NewZipf(rnd, 1.2, 1.1, uint64(conf.Clusters-1))
	globalTags := NewTagsDef(rnd, "global", "gv", conf.GlobalTags, conf.GlobalTags)
	metrics := generateMetrics(rnd, conf)
	if glog.V(1) {
		glog.Info("Generated ", len(metrics), " metrics")
	}
	machines := make([]*Machine, conf.Machines)
	for machine := 0; machine < conf.Machines; machine++ {
		var tags Tags
		machineName := genName("machine-", machine)
		clusterIdx := clusterGen.Uint64()
		tags = append(tags, globalTags.selectTags(conf.MinimumTags)...)
		tags = append(tags, clusters[clusterIdx].tags.selectTags(conf.MinimumTags)...)
		tags = append(tags, Tag{Name: "machine", Value: machineName})
		tags = append(tags, Tag{Name: "rack", Value: strconv.Itoa(machine / 48)})
		metricsCount := rnd.Intn(conf.MaxMetrics - conf.BaseMetrics)
		machines[machine] = NewMachine(machineName, tags,
			startTime+uint64(rnd.Intn(conf.StartSplay)), metricsCount+conf.BaseMetrics)
		if glog.V(1) {
			glog.Info("Machine ", machineName, " cluster=", clusterIdx,
				" tags=", tags, " startTime=", machines[machine].lastTs)
		}
		for metric := 0; metric < conf.BaseMetrics; metric++ {
			gen, name := metrics[metric].genMaker()
			selectedTags := metrics[metric].tags.selectTags(conf.TagsPerMetric)
			machines[machine].AddTimeseriesWithTags(
				metrics[metric].namespace,
				name,
				selectedTags,
				gen,
				metrics[metric].period)
		}
		if glog.V(1) {
			glog.Info("Machine ", machineName, " has ", conf.BaseMetrics, " base metrics")
		}
		metricsSelected := make(map[int]bool)
		for i := 0; i < metricsCount; i++ {
			metricsSelected[rnd.Intn(conf.MaxMetrics-conf.BaseMetrics)+conf.BaseMetrics] = true
		}
		if glog.V(2) {
			glog.Info("Machine ", machineName, " selected ", metricsSelected)
		}
		for metric := range metricsSelected {
			gen, name := metrics[metric].genMaker()
			selectedTags := metrics[metric].tags.selectTags(conf.TagsPerMetric)
			if glog.V(2) {
				glog.Info("Machine ", machineName, " adding metric ", metrics[metric])
			}
			machines[machine].AddTimeseriesWithTags(
				metrics[metric].namespace,
				name,
				selectedTags,
				gen,
				metrics[metric].period)
		}
		if glog.V(1) {
			glog.Info("Machine ", machineName, " has ", len(metricsSelected), " unique metrics")
		}
	}
	return &Simulator{conf, startTime, machines, globalTags}
}

// Run simulator for specified runFor time uints (whatever you use timestamp for, usually seconds)
// Callback will get tagged points usable to be sent to monitoring systems
func (s *Simulator) Run(shard int, shardCount int, runTo int64, cb func(points *[]TaggedPoints)) {
	s.currentTime = uint64(runTo)
	for i := shard; i < len(s.machines); i += shardCount {
		tick := s.machines[i].Tick(s.currentTime)
		if len(*tick) > 0 {
			cb(tick)
		}
	}
}

type tagsDef struct {
	tags              Tags
	tagDistribution   *rand.Zipf
	countDistribution *rand.Zipf
}

type metricDef struct {
	namespace string
	period    uint64
	tags      tagsDef
	genMaker  func() (generator.Generator, string)
}

func NewTagsDef(rnd *rand.Rand, namePrefix, valuePrefix string, maxTags int, maxCount int) tagsDef {
	tags := generateTags(rnd, namePrefix, valuePrefix, maxTags)
	if len(tags) == 0 {
		panic("Tags is zero for " + namePrefix)
	}
	return tagsDef{
		tags:              tags,
		tagDistribution:   rand.NewZipf(rnd, 1.3, 1.1, uint64(maxTags-1)),
		countDistribution: rand.NewZipf(rnd, 1.2, 1.1, uint64(maxCount-1)),
	}
}
func (td *tagsDef) selectTags(minimum int) Tags {
	if len(td.tags) == 0 {
		return nil
	}
	count := int(td.countDistribution.Uint64())
	if count < minimum {
		count = minimum
	}
	tags := make(Tags, 0, count)
	toGo := count
	for toGo > 0 {
		for idx := len(tags); idx < count; idx++ {
			tagIdx := int(td.tagDistribution.Uint64())
			tags = append(tags, td.tags[tagIdx])
		}
		sort.Sort(&tags)
		deduplicateTags(&tags)
		toGo -= len(tags)
	}
	return tags
}

func generateClusters(rnd *rand.Rand, conf Configuration) []clusterDef {
	clusters := make([]clusterDef, conf.Clusters)
	for cID := 0; cID < conf.Clusters; cID++ {
		clusters[cID].cID = cID
		clusters[cID].tags = NewTagsDef(rnd,
			"svc"+strconv.Itoa(cID)+"-", "v", conf.ClusterTags, conf.ClusterTags)
		if glog.V(2) {
			glog.Info("Cluster ", cID, " tags=", clusters[cID].tags.tags)
		}
	}
	return clusters
}

func genName(prefix string, id int) string {
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.WriteString(strconv.Itoa(id))
	return buf.String()
}

func genOrCache(cache *map[uint64]string, prefix string, key uint64) string {
	name, ok := (*cache)[key]
	if !ok {
		name = genName(prefix, int(key))
		(*cache)[key] = name
	}
	return name
}

func generateTags(rnd *rand.Rand, namePrefix, valuePrefix string, numTags int) Tags {
	if numTags == 0 {
		panic("numTags shouldn't be zero for " + namePrefix)
	}
	tagsCache := make(map[uint64]string)
	valuesCache := make(map[uint64]string)
	tagsGen := rand.NewZipf(rnd, 1.3, 1.1, uint64(numTags-1))
	valuesGen := rand.NewZipf(rnd, 1.3, 1.1, uint64(numTags-1))
	tags := make(Tags, 0, numTags)
	tagsToGo := numTags
	for tagsToGo > 0 {
		for tagIdx := len(tags); tagIdx < numTags; tagIdx++ {
			tagName := genOrCache(&tagsCache, namePrefix, tagsGen.Uint64())
			value := genOrCache(&valuesCache, valuePrefix, valuesGen.Uint64())
			tags = append(tags, Tag{tagName, value})
		}
		sort.Sort(&tags)
		deduplicateTags(&tags)
		tagsToGo = numTags - len(tags)
	}
	return tags
}

func deduplicateTags(tags *Tags) {
	var newTags Tags
	if len(*tags) < 2 {
		return
	}
	prev := 0
	newTags = append(newTags, (*tags)[prev])
	for i := 1; i < len(*tags); i++ {
		if (*tags)[i].Value != (*tags)[prev].Value {
			prev = i
			newTags = append(newTags, (*tags)[prev])
		}
	}
	*tags = newTags
}

func generateMetrics(rnd *rand.Rand, conf Configuration) []metricDef {
	count := conf.MaxMetrics
	periodDF := rand.NewZipf(rnd, 1.3, 1.2, 19)
	namespaces := count/conf.MetricsPerNamespace + 1
	metrics := make([]metricDef, count)
	for i := 0; i < count; i++ {
		genNum := rnd.Intn(10) // make randomwalk more frequent
		scale := float64(rnd.Intn(1e6)) * 0.001
		ns := rnd.Intn(namespaces)
		metrics[i].namespace = genName("ns", ns)
		metrics[i].period = uint64(15 * (int(periodDF.Uint64()) + 1))
		tagsCount := rnd.Intn(conf.TagsPerMetric * 2)
		if tagsCount > conf.TagsPerMetric/2 {
			metrics[i].tags = NewTagsDef(rnd, "mtag", "mv",
				conf.TagsPerMetric, tagsCount/2)
		}
		metricID := i
		metrics[i].genMaker = func() (generator.Generator, string) {
			genNum := genNum
			scale := scale
			var gen generator.Generator
			var metricPrefix string
			switch genNum {
			case 4:
				gen = generator.NewIncreasingGenerator(rnd, 0.8, 0.01, 0.1, 0.0, 100.0)
				metricPrefix = "metricI"
			case 1:
				gen = generator.NewSpikesGenerator(rnd, 100.0)
				metricPrefix = "metricS"
			case 3:
				gen = generator.NewCyclicGenerator(rnd, 0.0, 100.0)
				metricPrefix = "metricC"
			default:
				gen = generator.NewRandomWalkGenerator(rnd, 1.0, 0.0, 100)
				metricPrefix = "metricR"
			}
			return generator.NewScalingGenerator(gen, scale), genName(metricPrefix, metricID)
		}
	}
	return metrics
}
