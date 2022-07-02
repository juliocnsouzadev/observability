package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	appName = os.Getenv("APP_NAME")

	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "app_api",
			Name:      "counter",
			Help:      fmt.Sprintf("This is a counter for %s", appName),
		})

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "app_api",
			Name:      "gauge",
			Help:      fmt.Sprintf("This is a gauge for %s", appName),
		})

	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "app_api",
			Name:      "histogram",
			Help:      fmt.Sprintf("This is a histogram for %s", appName),
		})

	summary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "app_api",
			Name:      "summary",
			Help:      fmt.Sprintf("This is a summary for %s", appName),
		})
)

func main() {
	rand.Seed(time.Now().Unix())

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "prom_request_time",
		Help: "Time it has taken to retrieve the metrics",
	}, []string{"time"})

	prometheus.Register(histogramVec)

	http.Handle("/metrics", newHandlerWithHistogram(promhttp.Handler(), histogramVec))

	prometheus.MustRegister(counter)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(histogram)
	prometheus.MustRegister(summary)

	go func() {

		factor := float64(rand.Intn(10))
		
		for {
			counter.Add(rand.Float64() * factor)
			gauge.Add(rand.Float64()*factor - 5)
			histogram.Observe(rand.Float64() * factor)
			summary.Observe(rand.Float64() * factor)

			time.Sleep(time.Second * 2)
		}
	}()

	log.Fatal(http.ListenAndServe(":2112", nil))
}

func newHandlerWithHistogram(handler http.Handler, histogram *prometheus.HistogramVec) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		status := http.StatusOK

		defer func() {
			histogram.WithLabelValues(fmt.Sprintf("%d", status)).Observe(time.Since(start).Seconds())
		}()

		if req.Method == http.MethodGet {
			handler.ServeHTTP(w, req)
			return
		}
		status = http.StatusBadRequest

		w.WriteHeader(status)
	})
}
