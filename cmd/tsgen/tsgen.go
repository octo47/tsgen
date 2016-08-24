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
	"sync/atomic"
	"time"

	"github.com/octo47/tsgen"
)

var parallel = flag.Int("parallel", 1, "parallelize generation")
var seed = flag.Int64("seed", 1234, "seed randomization")
var machines = flag.Int("machines", 20, "machines run simulate")
var shard = flag.Int("shard", -1, "shard to run")
var clusters = flag.Int("clusters", 2, "number of clusters to generate")
var metrics = flag.Int("metrics", 1000, "metrics at total to simulate")
var pollPeriod = flag.Uint64("poll", 300, "simulated metrics publish period")
var startTimestamp = flag.Uint64("start", 0, "start from timestamp")
var duration = flag.Duration("duration", 30*time.Minute, "duration of generation")
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
	exit := int32(0)
	startTs := *startTimestamp
	if startTs == 0 {
		startTs = uint64(time.Now().Unix())
	}
	rnd := rand.New(rand.NewSource(*seed))
	conf := tsgen.NewConfiguration(*machines, *metrics, *clusters)
	sim := tsgen.NewSimulator(rnd, conf, startTs)
	if *nogen {
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// Block until a signal is received.
		<-c
		atomic.StoreInt32(&exit, 1)
	}()

	var wg sync.WaitGroup
	if *shard == -1 {
		for sh := 0; sh < *parallel; sh++ {
			wg.Add(1)
			go runShard(sim, sh, &wg, &exit)
		}
	} else {
		wg.Add(1)
		go runShard(sim, *shard, &wg, &exit)

	}
	wg.Wait()
}

func runShard(sim *tsgen.Simulator, shard int, wg *sync.WaitGroup, exit *int32) {
	defer wg.Done()
	for atomic.LoadInt32(exit) != 1 {
		sim.Run(shard, *parallel, *pollPeriod, func(tp *[]tsgen.TaggedPoints) {
			for i := range *tp {
				tagstr := (*tp)[i].Tags.FormatSeparated(' ')
				for _, point := range *(*tp)[i].Datapoints {
					if atomic.LoadInt32(exit) == 1 {
						return
					}
					var fullMetricName bytes.Buffer
					_, _ = fullMetricName.WriteString(*(*tp)[i].Namespace)
					_, _ = fullMetricName.WriteRune('.')
					_, _ = fullMetricName.WriteString(*(*tp)[i].MetricName)
					fmt.Println("put", fullMetricName.String(),
						point.Timestamp, point.Value, tagstr)
				}
			}
		})
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
