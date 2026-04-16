package healthcheck_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elizamthomas1994/vaultpipe/internal/healthcheck"
)

func TestCheck_ReachableActive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := healthcheck.New(srv.URL, 2*time.Second)
	res := c.Check(context.Background())

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if !res.Reachable {
		t.Error("expected reachable=true")
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}
	if res.Latency == 0 {
		t.Error("expected non-zero latency")
	}
}

func TestCheck_StandbyIsReachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
	}))
	defer srv.Close()

	c := healthcheck.New(srv.URL, 2*time.Second)
	res := c.Check(context.Background())

	if !res.Reachable {
		t.Error("expected standby vault to be reachable")
	}
}

func TestCheck_ServerError_NotReachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := healthcheck.New(srv.URL, 2*time.Second)
	res := c.Check(context.Background())

	if res.Reachable {
		t.Error("expected reachable=false for 500")
	}
}

func TestCheck_Unreachable(t *testing.T) {
	c := healthcheck.New("http://127.0.0.1:19999", 500*time.Millisecond)
	res := c.Check(context.Background())

	if res.Err == nil {
		t.Error("expected error for unreachable host")
	}
	if res.Reachable {
		t.Error("expected reachable=false")
	}
}

func TestCheck_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := healthcheck.New(srv.URL, 2*time.Second)
	res := c.Check(ctx)

	if res.Err == nil {
		t.Error("expected error for cancelled context")
	}
}
