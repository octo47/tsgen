package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/octo47/tsgen"
)

var parallel = flag.Int("parallel", 1, "parallelize generation")
var seed = flag.Int64("seed", 1234, "seed randomization")
var machines = flag.Int("machines", 20, "machines run simulate")
var metrics = flag.Int("metrics", 1000, "metrics at total to simulate")
var pollPeriod = flag.Uint64("poll", 300, "simulated metrics publish period")

func main() {
	var wg sync.WaitGroup
	exit := int32(0)
	rnd := rand.New(rand.NewSource(*seed))
	conf := tsgen.NewConfiguration(*machines, *metrics)
	sim := tsgen.NewSimulator(rnd, conf, uint64(time.Now().Unix()))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// Block until a signal is received.
		<-c
		atomic.StoreInt32(&exit, 1)
	}()
	for shard := 0; shard < *parallel; shard++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			shard := shard
			for atomic.LoadInt32(&exit) != 1 {
				sim.Run(shard, *parallel, *pollPeriod, func(tp *[]tsgen.TaggedPoints) {
					for i := range *tp {
						fmt.Println(i)
						//for _, point := range (*tp)[i].points {
						//	fmt.Println(point.Timestamp)
						//}
					}
				})
			}
		}()
	}
	wg.Wait()
}
