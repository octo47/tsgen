package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/octo47/tsgen/simulator"
)

var parallel = flag.Int("parallel", 1, "parallelize generation")
var seed = flag.Int64("seed", 1234, "seed randomization")
var machines = flag.Int("machines", 20, "machines run simulate")
var shard = flag.Int("shard", -1, "shard to run")
var clusters = flag.Int("clusters", 2, "number of clusters to generate")
var metrics = flag.Int("metrics", 1000, "metrics at total to simulate")
var pollPeriod = flag.Duration("poll", 1*time.Second, "simulated metrics publish period")
var startTimestamp = flag.Int64("start",
	(time.Now().Add(-5 * time.Minute)).Unix(), "start from timestamp")
var nogen = flag.Bool("nogen", false, "Do not run actual simluation, prepare only")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write mem profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	defer memoryProfile()
	exit := make(chan struct{})
	startTs := *startTimestamp
	rnd := rand.New(rand.NewSource(*seed))
	conf := simulator.NewConfiguration(*machines, *metrics, *clusters)
	sim := simulator.NewSimulator(rnd, conf, uint64(startTs))
	if *nogen {
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// Block until a signal is received.
		<-c
		close(exit)
	}()

	var wg sync.WaitGroup
	if *shard == -1 {
		for sh := 0; sh < *parallel; sh++ {
			wg.Add(1)
			go runShard(sim, sh, &wg, exit)
		}
	} else {
		wg.Add(1)
		go runShard(sim, *shard, &wg, exit)

	}
	wg.Wait()
}

func runShard(sim *simulator.Simulator, shard int, wg *sync.WaitGroup, exit chan struct{}) {
	defer wg.Done()

	currTime := time.Unix(*startTimestamp, 0)

	// catch up time first.
	for currTime.Before(time.Now()) {
		select {
		case <-exit:
			return
		default:
		}
		sim.Run(shard, *parallel, currTime.Unix(), consumePoints)
		currTime = currTime.Add(*pollPeriod)
	}

	ticker := time.NewTicker(*pollPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-exit:
			return
		case <-ticker.C:
			sim.Run(shard, *parallel, time.Now().Unix(), consumePoints)
		}
	}
}

func consumePoints(tp *[]simulator.TaggedPoints) {
	for i := range *tp {
		tagstr := (*tp)[i].Tags.FormatSeparated(' ')
		for _, point := range *(*tp)[i].Datapoints {
			var fullMetricName bytes.Buffer
			_, _ = fullMetricName.WriteString(*(*tp)[i].Namespace)
			_, _ = fullMetricName.WriteRune('.')
			_, _ = fullMetricName.WriteString(*(*tp)[i].MetricName)
			fmt.Println("put", fullMetricName.String(),
				point.Timestamp, point.Value, tagstr)
		}
	}
}

func memoryProfile() {
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.WriteHeapProfile(f)
		_ = f.Close()
		return
	}
}
