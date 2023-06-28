package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
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

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	outputDirPath := filepath.Join(*outputDir, time.Now().Format(time.RFC850))
	log.Infof("Creating dir for this benchmark at `%s`", outputDirPath)
	if err := os.MkdirAll(outputDirPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	reqBody, err := json.Marshal(RequestBody{
		queryStringParameters: map[string]int{
			"incrementLimit": 1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	endpoint := "<API_GATEWAY_ROUTE>"
	if snapStartEnabled {
		endpoint = endpoint + "hellojava_SnapStartEnabled"
	} else {
		endpoint = endpoint + "hellojava_SnapStartDisabled"
	}
	req, _ := http.NewRequest(http.MethodGet, endpoint, bytes.NewReader(reqBody))
	req = req.WithContext(ctx)
	signer := v4.NewSigner(cfg.Credentials)
	_, err = signer.Sign(req, nil, "execute-api", cfg.Region, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal(err)
		return
	}
	fmt.Println(res.Body)
}
