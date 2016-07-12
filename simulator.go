package tsgen

import (
	"bytes"
	"math/rand"
	"sort"
	"strconv"

	"github.com/golang/glog"
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
	UniqueTags int
	// minimun tags
	MinimumTags int
	// StartSplay configures how far machines start time could drift away from startTime
	StartSplay int
}

type clusterDef struct {
	cID  int
	tags tagsDef
}

func NewSimulator(rnd *rand.Rand, conf Configuration, startTime uint64) *Simulator {
	clusters := generateClusters(rnd, conf)
	clusterGen := rand.NewZipf(rnd, 1.2, 1.1, uint64(conf.Clusters-1))
	globalTags := NewTagsDef(rnd, conf.GlobalTags, conf.GlobalTags)
	machines := make([]*Machine, conf.Machines)
	for machine := 0; machine < conf.Machines; machine++ {
		var tags Tags
		machineName := genName("host-", machine)
		clusterIdx := clusterGen.Uint64()
		tags = append(tags, globalTags.selectTags(conf.MinimumTags)...)
		tags = append(tags, clusters[clusterIdx].tags.selectTags(conf.MinimumTags)...)
		machines[machine] = NewMachine(machineName, tags,
			startTime+uint64(rnd.Intn(conf.StartSplay)))
		if glog.V(1) {
			glog.Info("Machine ", machineName, " cluster=", clusterIdx,
				" tags=", tags)
		}
	}
	return &Simulator{conf, startTime, machines, globalTags}
}

// Run simulator for specified runFor time uints (whatever you use timestamp for, usually seconds)
// Callback will get tagged points usable to be sent to monitoring systems
func (s *Simulator) Run(shard int, shardCount int, runFor uint64, cb func(points *[]TaggedPoints)) {
	s.currentTime += runFor
	for i := 0; i < len(s.machines); i += shardCount {
		cb(s.machines[i+shard].Tick(s.currentTime))
	}
}

type tagsDef struct {
	tags              Tags
	tagDistribution   *rand.Zipf
	countDistribution *rand.Zipf
}

func NewTagsDef(rnd *rand.Rand, maxTags int, maxCount int) tagsDef {
	return tagsDef{
		tags:              generateTags(rnd, maxTags),
		tagDistribution:   rand.NewZipf(rnd, 1.2, 1.1, uint64(maxTags-1)),
		countDistribution: rand.NewZipf(rnd, 1.2, 1.1, uint64(maxCount-1)),
	}
}
func (td *tagsDef) selectTags(minimum int) Tags {
	count := int(td.countDistribution.Uint64())
	if count < minimum {
		count = minimum
	}
	tags := make(Tags, count)
	for idx := 0; idx < len(tags); idx++ {
		tagIdx := int(td.tagDistribution.Uint64())
		tags[idx] = td.tags[tagIdx]
	}
	return tags
}

func generateClusters(rnd *rand.Rand, conf Configuration) []clusterDef {
	clusters := make([]clusterDef, conf.Clusters)
	for cID := 0; cID < conf.Clusters; cID++ {
		clusters[cID].cID = cID
		clusters[cID].tags = NewTagsDef(rnd, conf.UniqueTags, conf.UniqueTags)
		if glog.V(1) {
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

func generateTags(rnd *rand.Rand, numTags int) Tags {
	tagsCache := make(map[uint64]string)
	valuesCache := make(map[uint64]string)
	tags := make(Tags, numTags)
	tagsGen := rand.NewZipf(rnd, 1.2, 1.1, uint64(numTags-1))
	valuesGen := rand.NewZipf(rnd, 1.1, 1.1, uint64(numTags-1))
	for tagIdx := 0; tagIdx < numTags; tagIdx++ {
		tagName := genOrCache(&tagsCache, "tag-", tagsGen.Uint64())
		value := genOrCache(&valuesCache, "value-", valuesGen.Uint64())
		tags[tagIdx] = Tag{tagName, value}
	}
	sort.Sort(&tags)
	deduplicateTags(&tags)
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
		if (*tags)[i] != (*tags)[prev] {
			prev = i
			newTags = append(newTags, (*tags)[prev])
		}
	}
	tags = &newTags
}
