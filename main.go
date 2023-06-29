package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"sync"
	"time"
)

type RequestBody struct {
	queryStringParameters map[string]int `json:"queryStringParameters"`
}

var outputDir = flag.String("o", "./latency_samples/", "Output directory for latency samples (default: ./latency_samples/)")
var burstCount = flag.Int("b", 1, "Number of bursts to send (default: 1)")
var burstSize = flag.Int("s", 1, "Number of requests to send in a burst (default: 1)")
var burstIAT = flag.Int("i", 10000, "Inter-arrival time between bursts in milliseconds (default: 10000)")
var parallel = flag.Bool("p", false, "Send requests in parallel (default: false)")
var snapStartEnabled = flag.Bool("snapstart", false, "Enable SnapStart (default: false)")

func main() {
	log.Info("Starting application...")
	flag.Parse()
	// Format current time in RFC3339 format
	currentTime := time.Now().Format(time.RFC3339)

	endpoint := "<API_GATEWAY_ROUTE>"
	//	if *snapStartEnabled {
	//		endpoint = endpoint + "hellojava_SnapStartEnabled"
	//	} else {
	//		endpoint = endpoint + "hellojava_SnapStartDisabled"
	//	}
	snapStartEnabledEndpoint := endpoint + "hellojava_SnapStartEnabled"
	snapStartDisabledEndpoint := endpoint + "hellojava_SnapStartDisabled"
	snapStartEnabledData := []int{}
	snapStartDisabledData := []int{}

	var wg sync.WaitGroup
	Log.Infof("Running benchmarks...")
	wg.Add(2)
	go RunBenchMark(&wg, snapStartEnabledEndpoint, *burstCount, true, &snapStartEnabledData)
	go RunBenchMark(&wg, snapStartDisabledEndpoint, *burstCount, false, &snapStartDisabledData)
	wg.Wait()
	Log.Infof("Benchmarks completed.")

	Log.Infof("Writing data to CSV file...")
	wg.Add(2)
	go WriteDataToFile(&wg, snapStartEnabledData, *outputDir+currentTime+"_snapstart_enabled.csv")
	go WriteDataToFile(&wg, snapStartDisabledData, *outputDir+currentTime+"_snapstart_disabled.csv")
	wg.Wait()
	Log.Infof("Written data to files %s and %s", *outputDir+currentTime+"_snapstart_enabled.csv", *outputDir+currentTime+"_snapstart_disabled.csv")
}

// RunCommandAndLog runs a command in the terminal, logs the result and returns it
func RunCommandAndLog(cmd *exec.Cmd) string {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s: %s", fmt.Sprint(err.Error()), stderr.String())
	}
	log.Debugf("Command result: %s", out.String())
	return out.String()
}

func RunBenchMark(wg *sync.WaitGroup, endpoint string, iterations int, snapStartEnabled bool, data *[]int) {
	defer wg.Done()
	log.Infof("Running benchmark with %d iterations, SnapStart enabled: %t", iterations, snapStartEnabled)
	for i := 0; i < iterations; i++ {
		log.Infof("Iteration %d", i)
		command := `curl -X GET -H "x-api-key: <API_KEY>" -H "Content-Type: application/json" ` + endpoint
		log.Info(command)
		startTime := time.Now()
		log.Infof("Sending request at %s", startTime)
		RunCommandAndLog(exec.Command("sh", "-c", command))
		endTime := time.Now()
		log.Infof("Request completed at %s", endTime)
		latency := endTime.Sub(startTime).Milliseconds()
		log.Infof("Time taken for request: %d", latency)
		*data = append(*data, int(latency))
		time.Sleep(burstIAT * time.Millisecond)
	}
}
