// sseload — M12 SSE load spike (plan μ46–μ49), Go-client stand-in for
// k6+xk6-sse (the extension is unavailable in this environment; see
// docs/testing/tooling-evaluation.md §2).
//
// Profiles against the demo binary's SSE endpoint DIRECTLY (never the
// Firebase-buffered canonical URL, plan guardrail 6 / #329):
//
//	sseload -profile steady   -url http://localhost:8093/api/datastar/stream -clients 10 -duration 30s
//	sseload -profile reconnect -url ... -attempts 10
//	sseload -profile drop     -url ... -clients 50
//
// Prints a JSON report on stdout: arrival-gap distribution (the server
// patches every 2s), continuity violations, keep-alive observations,
// reconnect time-to-first-event, and post-drop service health.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"
)

type gapStats struct {
	MeanMs     float64 `json:"meanMs"`
	P95Ms      float64 `json:"p95Ms"`
	MaxMs      float64 `json:"maxMs"`
	Violations int     `json:"violations"` // gaps exceeding 2× the expected cadence
}

type report struct {
	Profile          string    `json:"profile"`
	URL              string    `json:"url"`
	Clients          int       `json:"clients"`
	DurationSec      float64   `json:"durationSec"`
	TotalEvents      int       `json:"totalEvents"`
	TotalKeepalives  int       `json:"totalKeepalives"`
	Gaps             gapStats  `json:"arrivalGaps"`
	Errors           []string  `json:"errors"`
	ReconnectFirstMs []float64 `json:"reconnectToFirstEventMs,omitempty"`
	PostDropHealthOK bool      `json:"postDropHealthOk,omitempty"`
}

func stream(url string, out chan<- time.Time, keepalives chan<- struct{}, errs chan<- error, stop <-chan struct{}) {
	go func() {
		defer close(out)

		client := &http.Client{Timeout: 0}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			errs <- err

			return
		}

		req.Header.Set("Accept", "text/event-stream")

		resp, err := client.Do(req)
		if err != nil {
			errs <- err

			return
		}
		defer resp.Body.Close()

		lines := make(chan string, 256)
		go func() {
			defer close(lines)

			sc := bufio.NewScanner(resp.Body)
			sc.Buffer(make([]byte, 64*1024), 1024*1024)

			for sc.Scan() {
				lines <- sc.Text()
			}
		}()

		inEvent := false

		for {
			select {
			case <-stop:
				return
			case line, ok := <-lines:
				if !ok {
					return
				}

				switch {
				case line == "":
					inEvent = false
				case len(line) > 2 && line[0] == ':':
					select {
					case keepalives <- struct{}{}:
					default:
					}
				case len(line) > 6 && line[:7] == "event: ":
					if !inEvent {
						inEvent = true

						out <- time.Now()
					}
				}
			}
		}
	}()
}

func summarize(gaps []time.Duration, cadence time.Duration) gapStats {
	if len(gaps) == 0 {
		return gapStats{}
	}

	sorted := append([]time.Duration(nil), gaps...)
	slices.Sort(sorted)

	var total time.Duration

	violations := 0

	for _, g := range sorted {
		total += g
		if g > 2*cadence {
			violations++
		}
	}

	p95 := sorted[int(float64(len(sorted)-1)*0.95)]

	return gapStats{
		MeanMs:     float64(total.Milliseconds()) / float64(len(sorted)),
		P95Ms:      float64(p95.Milliseconds()),
		MaxMs:      float64(sorted[len(sorted)-1].Milliseconds()),
		Violations: violations,
	}
}

func main() {
	url := flag.String(
		"url",
		"http://localhost:8093/api/datastar/stream",
		"SSE endpoint (demo binary, NOT a proxied URL)",
	)
	profile := flag.String("profile", "steady", "steady | reconnect | drop")
	clients := flag.Int("clients", 10, "concurrent clients")
	duration := flag.Duration("duration", 30*time.Second, "steady-state hold")
	attempts := flag.Int("attempts", 10, "reconnect attempts")
	cadence := flag.Duration("cadence", 2*time.Second, "expected event cadence")

	flag.Parse()

	rep := report{Profile: *profile, URL: *url, Clients: *clients, DurationSec: duration.Seconds()}

	switch *profile {
	case "steady":
		var (
			mu      sync.Mutex
			allGaps []time.Duration
		)

		events, keepalives, errs := make(chan time.Time), make(chan struct{}, 1024), make(chan error, *clients)

		stop := make(chan struct{})
		for range *clients {
			stream(*url, events, keepalives, errs, stop)
		}

		deadline := time.Now().Add(*duration)
		last := make(map[int]time.Time)
		idx := 0

		for time.Now().Before(deadline) {
			select {
			case _, ok := <-events:
				if !ok {
					continue
				}

				rep.TotalEvents++

				mu.Lock()

				gapIdx := idx % *clients
				if !last[gapIdx].IsZero() {
					allGaps = append(allGaps, time.Since(last[gapIdx]))
				}

				last[gapIdx] = time.Now()
				mu.Unlock()

				idx++
			case <-keepalives:
				rep.TotalKeepalives++
			case err := <-errs:
				rep.Errors = append(rep.Errors, err.Error())
			}
		}

		close(stop)

		rep.Gaps = summarize(allGaps, *cadence)

	case "reconnect":
		for range *attempts {
			start := time.Now()
			events, _, errs := make(chan time.Time), make(chan struct{}, 1024), make(chan error, 1)
			stop := make(chan struct{})
			stream(*url, events, keepaliveDrain(), errs, stop)

			select {
			case <-events:
				rep.ReconnectFirstMs = append(rep.ReconnectFirstMs, float64(time.Since(start).Milliseconds()))
				rep.TotalEvents++
			case err := <-errs:
				rep.Errors = append(rep.Errors, err.Error())
			case <-time.After(10 * time.Second):
				rep.Errors = append(rep.Errors, "timeout waiting for first event")
			}

			close(stop)
			time.Sleep(250 * time.Millisecond)
		}

	case "drop":
		var wg sync.WaitGroup

		stops := make([]chan struct{}, *clients)
		for i := range *clients {
			stops[i] = make(chan struct{})
			events, _, errs := make(chan time.Time), make(chan struct{}, 1024), make(chan error, 1)
			stream(*url, events, keepaliveDrain(), errs, stops[i])

			wg.Go(func() { ; <-events; rep.TotalEvents++ })
		}

		time.Sleep(3 * time.Second) // let all streams deliver their first patch

		for _, s := range stops { // synchronized client drop
			close(s)
		}

		wg.Wait()

		// Post-drop health: a fresh client must still receive events promptly.
		start := time.Now()
		events, _, errs := make(chan time.Time), make(chan struct{}, 1024), make(chan error, 1)
		stop := make(chan struct{})
		stream(*url, events, keepaliveDrain(), errs, stop)

		select {
		case <-events:
			rep.PostDropHealthOK = true
			rep.ReconnectFirstMs = append(rep.ReconnectFirstMs, float64(time.Since(start).Milliseconds()))
		case <-time.After(10 * time.Second):
			rep.PostDropHealthOK = false
		}

		close(stop)
	}

	out, _ := json.MarshalIndent(rep, "", "  ")
	fmt.Println(string(out))
	os.Exit(0)
}

func keepaliveDrain() chan struct{} {
	return make(chan struct{}, 1024)
}
