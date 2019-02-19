package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cha87de/tsprofiler/api"
	"github.com/cha87de/tsprofiler/models"
	tsprofiler "github.com/cha87de/tsprofiler/profiler"
	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
)

var byteCounter uint64
var counterAccess *sync.Mutex
var lastProfilerRun time.Time
var profiler api.TSProfiler

func init() {
	frequency := 30 // seconds
	counterAccess = &sync.Mutex{}
	profiler = tsprofiler.NewProfiler(models.Settings{
		Name:           "consumer_tsprofiler",
		BufferSize:     frequency * 20, // 10 Minutes
		States:         4,
		FilterStdDevs:  -1,
		History:        1,
		FixBound:       false,
		PeriodSize:     []int{6, 72, 144}, // 1h, 12h, 24h
		OutputFreq:     time.Duration(60) * time.Minute,
		OutputCallback: profileOutput,
	})
	go profileDumper(frequency)
}

func profileFlow(flow *flow.FlowMessage) {
	// simply increase counter
	counterAccess.Lock()
	byteCounter += flow.GetBytes()
	counterAccess.Unlock()
}

func profileDumper(frequency int) {
	// run periodically profileDump, e.g. every frequency seconds
	for {
		start := time.Now()
		profileDump()
		nextRun := start.Add(time.Duration(frequency) * time.Second)
		time.Sleep(nextRun.Sub(time.Now()))
	}
}

func profileDump() {
	// push counter to profiler && reset counter to zero
	counterAccess.Lock()
	timeDiff := time.Now().Sub(lastProfilerRun)
	byteRate := float64(byteCounter) / timeDiff.Seconds()
	max := float64(2)                      // Gb/s
	max = max * 1024 / 8                   // MB/s
	val := float64(byteRate / 1024 / 1024) // MB/s

	fmt.Printf("byteCounter: %d, rate: %.1f MB/s (max: %.1f MB/s)\n", byteCounter, val, max)

	// prepare metrics for profiler
	metrics := make([]models.TSInputMetric, 0)
	metrics = append(metrics, models.TSInputMetric{
		Name:     "net",
		Value:    float64(val),
		FixedMax: float64(max),
	})

	// send measurement to profiler
	tsinput := models.TSInput{
		Metrics: metrics,
	}
	profiler.Put(tsinput)

	lastProfilerRun = time.Now()
	byteCounter = 0
	counterAccess.Unlock()
}

func profileOutput(data models.TSProfile) {
	json, _ := json.Marshal(data)
	fmt.Printf("%s\n", json)
}
