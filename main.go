package main

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
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

	endpoint := "<API_GATEWAY_ROUTE>"
	if *snapStartEnabled {
		endpoint = endpoint + "hellojava_SnapStartEnabled"
	} else {
		endpoint = endpoint + "hellojava_SnapStartDisabled"
	}
	command := `curl -X GET -H "x-api-key: <API_KEY>" -H "Content-Type: application/json" ` + endpoint
	log.Info(command)
	startTime := time.Now()
	log.Infof("Sending request at %s", startTime)
	RunCommandAndLog(exec.Command("sh", "-c", command))
	endTime := time.Now()
	log.Infof("Request completed at %s", endTime)
	log.Infof("Time taken for request: %d", endTime.Sub(startTime).Milliseconds())
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
