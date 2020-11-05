package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

var gen token.Generator
var clusterID string
var psk string

func handler(w http.ResponseWriter, r *http.Request) {
	var tok token.Token
	var err error

	values := r.URL.Query()
	if values.Get("psk") != psk {
		http.Error(w, "wrong psk", http.StatusForbidden)
		return
	}

	tok, err = gen.Get(clusterID)
	if err != nil {
		http.Error(w, "failed to retrieve token", http.StatusServiceUnavailable)
	}
	log.Printf("Got token %v", gen.FormatJSON(tok))
	fmt.Fprintf(w, "%v\n", gen.FormatJSON(tok))
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
	log.Info("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
