package chiprometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Test_Logger(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	m := NewMiddleware("test")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/users/JoeBob", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequest("GET", "http://localhost:3000/users/Misty", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)
	body := recorder.Body.String()
	if !strings.Contains(body, reqsName) {
		t.Errorf("body does not contain request total entry '%s'", reqsName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", latencyName)
	}

	req1Count := `chi_request_duration_milliseconds_count{code="OK",method="GET",path="/ok",service="test"} 1`
	req2Count := `chi_request_duration_milliseconds_count{code="OK",method="GET",path="/users/JoeBob",service="test"} 1`
	req3Count := `chi_request_duration_milliseconds_count{code="OK",method="GET",path="/users/Misty",service="test"} 1`

	if !strings.Contains(body, req1Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req1Count)
	}
	if !strings.Contains(body, req2Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req2Count)
	}
	if !strings.Contains(body, req3Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req3Count)
	}
}

func Test_PatternLogger(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	m := NewPatternMiddleware("patternOnlyTest")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/users/JoeBob", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequest("GET", "http://localhost:3000/users/Misty", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)

	body := recorder.Body.String()

	if !strings.Contains(body, patternReqsName) {
		t.Errorf("body does not contain request total entry '%s'", patternReqsName)
	}
	if !strings.Contains(body, patternLatencyName) {
		t.Errorf("body does not contain request duration entry '%s'", patternLatencyName)
	}

	req1Count := `chi_pattern_request_duration_milliseconds_count{code="OK",method="GET",path="/ok",service="patternOnlyTest"} 1`
	joeBobCount := `chi_pattern_request_duration_milliseconds_count{code="OK",method="GET",path="/users/JoeBob",service="patternOnlyTest"} 1`
	mistyCount := `chi_pattern_request_duration_milliseconds_count{code="OK",method="GET",path="/users/Misty",service="patternOnlyTest"} 1`
	firstNamePatternCount := `chi_pattern_request_duration_milliseconds_count{code="OK",method="GET",path="/users/{firstName}",service="patternOnlyTest"} 2`

	if !strings.Contains(body, req1Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req1Count)
	}
	if strings.Contains(body, joeBobCount) {
		t.Errorf("body should not contain Joe Bob count summary '%s'", joeBobCount)
	}
	if strings.Contains(body, mistyCount) {
		t.Errorf("body should not contain Misty count summary '%s'", mistyCount)
	}
	if !strings.Contains(body, firstNamePatternCount) {
		t.Errorf("body does not contain first name pattern count summary '%s'", firstNamePatternCount)
	}
}

func Test_MultipleLoggers(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	mid := NewMiddleware("pathTest")
	m := NewPatternMiddleware("patternTest")

	n.Use(mid)
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/users/JoeBob", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequest("GET", "http://localhost:3000/users/Misty", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)

	body := recorder.Body.String()

	if !strings.Contains(body, patternReqsName) {
		t.Errorf("body does not contain request total entry '%s'", patternReqsName)
	}
	if !strings.Contains(body, patternLatencyName) {
		t.Errorf("body does not contain request duration entry '%s'", patternLatencyName)
	}

	req1Count := `chi_pattern_request_duration_milliseconds_count{code="OK",method="GET",path="/ok",service="patternTest"} 1`
	joeBobCount := `chi_request_duration_milliseconds_count{code="OK",method="GET",path="/users/JoeBob",service="pathTest"} 1`
	mistyCount := `chi_request_duration_milliseconds_count{code="OK",method="GET",path="/users/Misty",service="pathTest"} 1`
	firstNamePatternCount := `chi_pattern_request_duration_milliseconds_count{code="OK",method="GET",path="/users/{firstName}",service="patternTest"} 2`

	if !strings.Contains(body, req1Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req1Count)
	}
	if !strings.Contains(body, joeBobCount) {
		t.Errorf("body does not contain Joe Bob count summary '%s'", joeBobCount)
	}
	if !strings.Contains(body, mistyCount) {
		t.Errorf("body does not contain Misty count summary '%s'", mistyCount)
	}
	if !strings.Contains(body, firstNamePatternCount) {
		t.Errorf("body does not contain first name pattern count summary '%s'", firstNamePatternCount)
	}
}
