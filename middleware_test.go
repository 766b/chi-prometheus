package chiprometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
)

func Test_Logger(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	m := NewMiddleware("test")
	n.Use(m)

	n.Handle("/metrics", prometheus.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	body := recorder.Body.String()
	if !strings.Contains(body, "test_requests_total") {
		t.Errorf("body does not contain request total entry '%s'", "test_requests_total")
	}
	if !strings.Contains(body, "test_duration_milliseconds") {
		t.Errorf("body does not contain request duration entry '%s'", "test_duration_milliseconds")
	}
}
