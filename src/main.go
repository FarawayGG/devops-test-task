package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

type server struct {
	redis  redis.UniversalClient
	logger *zap.Logger
}

var indexHandlerCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "indexHandler_request_count",
		Help: "No of request handled by indexHandler handler",
	},
	[]string{"status"}, // labels
)

var indexHandlerLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_index_handler_duration_seconds",
		Help:    "Latency of index_handler request in second.",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
	},
	[]string{"status"},
)

func main() {
	prometheus.MustRegister(indexHandlerCounter)
	prometheus.MustRegister(indexHandlerLatency)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("unable to initialize logger")
	}

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{os.Getenv("REDIS_ADDR")},
		Password: "",
		DB:       0,
	})

	srv := &server{
		redis:  rdb,
		logger: logger,
	}

	router := httprouter.New()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info("Metric server started on port 8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	router.GET("/", srv.indexHandler)
	logger.Info("server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var status string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		indexHandlerLatency.WithLabelValues(status).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
		indexHandlerCounter.WithLabelValues(status).Inc()
	}()

	var v string
	var err error
	if v, err = s.redis.Get(context.Background(), "updated_time").Result(); err != nil {
		s.logger.Info("updated_time not found, setting it")
		v = time.Now().Format("2006-01-02 03:04:05")
		s.redis.Set(context.Background(), "updated_time", v, 5*time.Second)
	} else {
		s.logger.Info("got updated_time")
	}

	if err != nil {
		status = "success"
	} else {
		status = "error"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello world: updated_time=%s\n", v)
}
