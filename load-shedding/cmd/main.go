package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type LoadShedder struct {
	isOverLoaded atomic.Bool
}

func NewLoadShedder(ctx context.Context, checkInterval, overloadFactor time.Duration) *LoadShedder {
	ls := LoadShedder{}

	go ls.runOverLoadDetector(ctx, checkInterval, overloadFactor)

	return &ls
}

func (ls *LoadShedder) runOverLoadDetector(ctx context.Context, checkInterval, overloadFactor time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	// Start with a fresh start time
	startTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// check how long it took to process the latest batch of requests
			elapsed := time.Since(startTime)
			if elapsed > overloadFactor {
				// if it took longer than the overload factor, we're overloaded
				ls.isOverLoaded.Store(true)
			} else {
				// Otherwise, we're not overloaded
				ls.isOverLoaded.Store(false)
			}
			// Reset the start time
			startTime = time.Now()
		}
	}
}

func (ls *LoadShedder) IsOverLoaded() bool {
	return ls.isOverLoaded.Load()
}

type Handler struct {
	ls *LoadShedder
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if h.ls.IsOverLoaded() {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, http.StatusText(http.StatusServiceUnavailable))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
func NewHandler(ls *LoadShedder) *Handler {
	return &Handler{ls: ls}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The load shedder will check every 100ms if the last batch of requests took longer than 200ms
	ls := NewLoadShedder(ctx, 100*time.Millisecond, 200*time.Millisecond)

	h := NewHandler(ls)
	http.HandleFunc("GET /", h.GetHandler)
	log.Println("load-shedding listen on 8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
