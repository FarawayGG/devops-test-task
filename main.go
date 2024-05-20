package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type server struct {
	redis  redis.UniversalClient
	logger *zap.Logger
}

func main() {
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

	router.GET("/", srv.indexHandler)
	router.GET("/healthz", srv.healthzHandler)

	logger.Info("server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	metricsRouter := httprouter.New()
	metricsRouter.GET("/", Metrics(promhttp.Handler()))
	log.Fatal(http.ListenAndServe(":8081", metricsRouter))
}

func Metrics(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var v string
	var err error
	if v, err = s.redis.Get(context.Background(), "updated_time").Result(); err != nil {
		s.logger.Info("updated_time not found, setting it")
		v = time.Now().Format("2006-01-02 03:04:05")
		s.redis.Set(context.Background(), "updated_time", v, 5*time.Second)
	} else {
		s.logger.Info("got updated_time")
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello world: updated_time=%s\n", v)
}

func (s *server) healthzHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	healthReport := map[string]string{
		"redis": "healthy",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.redis.Ping(ctx).Err()
	if err != nil {
		s.logger.Error("redis ping failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unhealthy\n")
		return
	}

	if err != nil {
		healthReport["redis"] = "unhealthy"
	}

	jsonReport, err := json.Marshal(healthReport)
	if err != nil {
		s.logger.Error("failed to marshal health report", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to generate health report\n")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonReport)

}
