package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	healthcheckEndpointCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "custom",
			Name:      "http_request_healthcheck_count",
			Help:      "The total number of requests made to healthcheck endpoint",
		},
		[]string{"status"},
	)

	healthcheckEndpointLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "custom",
			Name:      "http_request_healthcheck_duration_seconds",
			Help:      "Latency of healthcheck endpoint requests in seconds",
		},
		[]string{"status"},
	)
)

func init() {
	// must register counter on init
	prometheus.MustRegister(healthcheckEndpointCounter)
	prometheus.MustRegister(healthcheckEndpointLatency)
}

func buildJSONResponse(statusCode int, message string) ([]byte, error) {
	var responseHTTP = make(map[string]interface{})

	responseHTTP["statusCode"] = statusCode
	responseHTTP["data"] = message

	response, err := json.Marshal(responseHTTP)
	if err != nil {
		return nil, err
	}

	return []byte(string(response)), nil
}

func returnHTTPResponse(statusCode int, message string) http.HandlerFunc {
	return func(writter http.ResponseWriter, req *http.Request) {
		responseJSONBytes, _ := buildJSONResponse(statusCode, message)

		var status string
		// counter for prometheus
		defer func() {
			healthcheckEndpointCounter.WithLabelValues(status).Inc()
		}()
		// latency timer for prometheus
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(_time float64) {
			healthcheckEndpointLatency.WithLabelValues(status).Observe(_time)
		}))
		defer func() {
			timer.ObserveDuration()
		}()

		// sleep from 0 to 1 secs randomly
		duration := (rand.Intn(1-0) + 0)
		time.Sleep(time.Duration(duration) * time.Second)

		writter.Header().Set("Content-Type", "application/json")
		writter.WriteHeader(statusCode)

		status = "success"

		writter.Write(responseJSONBytes)
	}
}

func ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func serve(address string, message string) {
	http.HandleFunc("/", returnHTTPResponse(http.StatusOK, message))
	http.HandleFunc("/healthcheck", returnHTTPResponse(http.StatusOK, "OK"))
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("[INFO] Listening on %s\n", address)
	http.ListenAndServe(address, nil)
}

func main() {
	port := ternary(os.Getenv("API_PORT") != "", os.Getenv("API_PORT"), "9000").(string) // default 9000
	address := fmt.Sprintf(":%s", port)
	fmt.Printf("[INFO] Address :: %s\n", address)

	message := ternary(os.Getenv("MESSAGE") != "", os.Getenv("MESSAGE"), "Hello World").(string) // default Hello World
	fmt.Printf("[INFO] Message :: %s\n", message)

	serve(address, message)
}
