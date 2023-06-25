package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"time"
)

var outputDir = flag.String("o", "./latency_samples/", "Output directory for latency samples (default: ./latency_samples/)")
var burstCount = flag.Int("b", 1, "Number of bursts to send (default: 1)")
var burstSize = flag.Int("s", 1, "Number of requests to send in a burst (default: 1)")
var burstIAT = flag.Int("i", 10000, "Inter-arrival time between bursts in milliseconds (default: 10000)")
var parallel = flag.Bool("p", false, "Send requests in parallel (default: false)")

func main() {
	log.Info("Starting application...")

}
