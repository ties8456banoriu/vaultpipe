// Package telemetry provides lightweight operation timing and metrics
// collection for vaultpipe commands and internal operations.
package telemetry

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Span represents a single timed operation.
type Span struct {
	Name     string
	Start    time.Time
	End      time.Time
	Err      error
	Duration time.Duration
}

// Finished returns true if the span has been ended.
func (s *Span) Finished() bool {
	return !s.End.IsZero()
}

// Collector records spans and can emit a summary.
type Collector struct {
	mu    sync.Mutex
	spans []*Span
	clock func() time.Time
}

// New returns a new Collector.
func New() *Collector {
	return &Collector{clock: time.Now}
}

// Start begins a new named span and returns it.
// Call span.Stop() or span.StopWithError() when the operation completes.
func (c *Collector) Start(name string) *Span {
	s := &Span{Name: name, Start: c.clock()}
	c.mu.Lock()
	c.spans = append(c.spans, s)
	c.mu.Unlock()
	return s
}

// Stop ends the span without an error.
func (c *Collector) Stop(s *Span) {
	c.stopSpan(s, nil)
}

// StopWithError ends the span and records the associated error.
func (c *Collector) StopWithError(s *Span, err error) {
	c.stopSpan(s, err)
}

func (c *Collector) stopSpan(s *Span, err error) {
	s.End = c.clock()
	s.Duration = s.End.Sub(s.Start)
	s.Err = err
}

// Spans returns a copy of all recorded spans.
func (c *Collector) Spans() []Span {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Span, len(c.spans))
	for i, s := range c.spans {
		out[i] = *s
	}
	return out
}

// WriteSummary writes a human-readable timing summary to w.
func (c *Collector) WriteSummary(w io.Writer) {
	spans := c.Spans()
	for _, s := range spans {
		status := "ok"
		if s.Err != nil {
			status = fmt.Sprintf("err: %s", s.Err)
		}
		fmt.Fprintf(w, "  %-30s %8.2fms  [%s]\n", s.Name, float64(s.Duration.Microseconds())/1000.0, status)
	}
}
