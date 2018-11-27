package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/edjumacator/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	n := chi.NewRouter()
	m := chiprometheus.NewPatternMiddleware("test_service")

	n.Use(m)

	n.Handle("/metrics", prometheus.Handler())
	n.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "slept %d milliseconds\n", sleep)
	})
	n.Get("/users/{firstName}", func(w http.ResponseWriter, r *http.Request) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "slept %d milliseconds\n", sleep)
	})
	n.Get("/users/{id}/contacts", func(w http.ResponseWriter, r *http.Request) {
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
