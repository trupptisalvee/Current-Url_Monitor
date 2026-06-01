package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ashish-dhakane/curlsc/internal/models"
)

func TestChecker_Check(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := &Checker{
		Timeout:   2 * time.Second,
		Method:    "GET",
		UserAgent: "test",
	}

	res := c.Check(ts.URL)

	if res.Status != models.StatusUp {
		t.Errorf("Expected status UP, got %s", res.Status)
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
}

func TestChecker_Check_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := &Checker{
		Timeout:   10 * time.Millisecond,
		Method:    "GET",
		UserAgent: "test",
	}

	res := c.Check(ts.URL)

	if res.Status != models.StatusError {
		t.Errorf("Expected status ERROR, got %s", res.Status)
	}
	if res.ErrorType != "TIMEOUT" {
		t.Errorf("Expected error type TIMEOUT, got %s", res.ErrorType)
	}
}
