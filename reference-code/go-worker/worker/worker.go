// Package worker provides a ticker-based background worker that processes
// tasks at regular intervals.
package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Worker handles periodic task processing with configurable tick rate
type Worker struct {
	tickRate  time.Duration
	lastTick  time.Time
	tickCount uint64
	running   bool
	mu        sync.RWMutex
}

// Stats represents the current state of the worker
type Stats struct {
	TickRate  string    `json:"tick_rate"`
	LastTick  time.Time `json:"last_tick"`
	TickCount uint64    `json:"tick_count"`
	Running   bool      `json:"running"`
}

// New creates a new Worker with the specified tick rate
func New(tickRate time.Duration) *Worker {
	return &Worker{
		tickRate: tickRate,
		running:  false,
	}
}

// Start begins the worker loop, processing tasks at each tick interval.
// It runs until the provided context is cancelled.
func (w *Worker) Start(ctx context.Context) {
	w.mu.Lock()
	w.running = true
	w.mu.Unlock()

	ticker := time.NewTicker(w.tickRate)
	defer ticker.Stop()

	slog.Info("Worker started", "tick_rate", w.tickRate)

	for {
		select {
		case <-ctx.Done():
			w.mu.Lock()
			w.running = false
			w.mu.Unlock()
			slog.Info("Worker stopped")
			return
		case <-ticker.C:
			w.process()
		}
	}
}

// process executes the work for a single tick
func (w *Worker) process() {
	w.mu.Lock()
	w.tickCount++
	currentTick := w.tickCount
	w.lastTick = time.Now()
	w.mu.Unlock()

	slog.Info("Worker tick processing", "tick", currentTick)
}

// GetStats returns the current worker statistics in a thread-safe manner
func (w *Worker) GetStats() Stats {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return Stats{
		TickRate:  w.tickRate.String(),
		LastTick:  w.lastTick,
		TickCount: w.tickCount,
		Running:   w.running,
	}
}
