package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type metric struct {
	Help  string
	Type  string
	Value int
}

var gen token.Generator
var clusterID string
var psk string
var wrongPskSleep = time.Second
var version = "undefined"
var metrics = map[string]*metric{
	"aws_iam_authenticator_proxy:tokens:total_requested": &metric{
		"Total number of token requested",
		"counter",
		0,
	},
	"aws_iam_authenticator_proxy:tokens:total_delivered": &metric{
		"Total number of token delivered",
		"counter",
		0,
	},
	"aws_iam_authenticator_proxy:tokens:total_errors": &metric{
		"Total number of token errored",
		"counter",
		0,
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	metrics["aws_iam_authenticator_proxy:tokens:total_requested"].Value += 1
	var tok token.Token
	var err error

	values := r.URL.Query()
	if values.Get("psk") != psk {
		metrics["aws_iam_authenticator_proxy:tokens:total_errors"].Value += 1
		time.Sleep(wrongPskSleep)
		http.Error(w, "wrong psk", http.StatusForbidden)
		return
	}

	tok, err = gen.Get(clusterID)
	if err != nil {
		metrics["aws_iam_authenticator_proxy:tokens:total_errors"].Value += 1
		http.Error(w, "failed to retrieve token", http.StatusServiceUnavailable)
		return
	}
	metrics["aws_iam_authenticator_proxy:tokens:total_delivered"].Value += 1
	log.Printf("Got token %v", gen.FormatJSON(tok))
	fmt.Fprintf(w, "%v\n", gen.FormatJSON(tok))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Everything running smooth")
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	for k, m := range metrics {
		fmt.Fprintf(w, "# HELP %s %s\n", k, m.Help)
		fmt.Fprintf(w, "# TYPE %s %s\n", k, m.Type)
		fmt.Fprintf(w, "%s{} %d\n", k, m.Value)
	}
}

func init() {
	var err error
	gen, err = token.NewGenerator(false, false)
	if err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}

	clusterID = os.Getenv("EKS_CLUSTER_ID")
	if clusterID == "" {
		log.Fatal("EKS_CLUSTER_ID must be set")
	}

	psk = os.Getenv("PSK")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)
	log.Infof("aws-iam-authenticator-proxy %s starting on port 8080", version)
	http.ListenAndServe(":8080", nil)
}
