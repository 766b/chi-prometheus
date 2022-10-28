package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	n := chi.NewRouter()
	m := chiprometheus.NewMiddleware("serviceName")

	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "slept %d milliseconds\n", sleep)
	})
	n.Get("/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not found")
	})

	http.ListenAndServe(":3000", n)
}
