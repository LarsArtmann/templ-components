package visualtest

import (
	"bytes"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// Demo server constants: health polling cadence, startup budget, and the
// module-path of the demo binary's package (built on demand from the sibling
// checkout the visualtest module replaces).
const (
	demoHealthPollInterval = 100 * time.Millisecond
	demoStartTimeout       = 60 * time.Second
	demoPackagePath        = "github.com/larsartmann/templ-components/examples/demo"
)

// DemoServer is a live examples/demo process bound to an ephemeral port, built
// and started on demand for browser-level tests (a11y sweep, click-through
// flows). Tests get the real demo routes — index showcase, forms, recipes,
// users — not a re-implementation.
type DemoServer struct {
	baseURL string
	cmd     *exec.Cmd
	// Log accumulates the server's combined stdout+stderr so tests (and the
	// CI smoke) can assert absence of 500 responses.
	Log bytes.Buffer
}

// BaseURL returns the demo server's root URL, e.g. http://127.0.0.1:41234.
func (s *DemoServer) BaseURL() string { return s.baseURL }

// StartDemoServer builds the examples/demo binary (cached by go's build cache,
// so the cost is paid once per machine, not per run) and starts it on an
// ephemeral port. Skips in -short mode. The process is killed on test cleanup.
func StartDemoServer(t *testing.T) *DemoServer {
	t.Helper()

	if testing.Short() {
		t.Skip("demo server e2e: skipped in -short mode")
	}

	binary := buildDemoBinary(t)
	port := reserveFreePort(t)

	server := &DemoServer{
		baseURL: "http://127.0.0.1:" + strconv.Itoa(port),
		cmd:     exec.Command(binary), //nolint:gosec,noctx // test fixture: locally built binary, lifecycle owned by the test
	}
	server.cmd.Stdout = &server.Log
	server.cmd.Stderr = &server.Log
	server.cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))

	if err := server.cmd.Start(); err != nil {
		t.Fatalf("visualtest[demo]: start server: %v", err)
	}

	done := make(chan error, 1)

	go func() { done <- server.cmd.Wait() }()

	waitForDemoHealth(t, server.baseURL)

	t.Cleanup(func() {
		if server.cmd.Process != nil {
			_ = server.cmd.Process.Kill()
		}
		<-done
	})

	return server
}

// buildDemoBinary compiles the demo into the test's temp dir and returns the
// binary path. GOEXPERIMENT=jsonv2 is forced because the demo links the
// errorpage package (encoding/json/v2 under Go 1.26's experiment).
func buildDemoBinary(t *testing.T) string {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "tc-demo")

	build := exec.Command("go", "build", "-o", binary, demoPackagePath) //nolint:gosec,noctx // test fixture: fixed package path, no user input
	build.Dir = "."
	build.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2", "GOWORK=off")

	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("visualtest[demo]: build demo binary: %v\n%s", err, output)
	}

	return binary
}

// reserveFreePort asks the kernel for a free TCP port and releases it for the
// demo server to bind. The classic bind-race window is accepted: tests pick
// from ~28k ephemeral ports while at most a handful of servers run.
func reserveFreePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("visualtest[demo]: reserve port: %v", err)
	}

	defer listener.Close()

	port, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("visualtest[demo]: reserved address is not TCP: %v", listener.Addr())
	}

	return port.Port
}

// waitForDemoHealth polls /health until the server responds 200, failing the
// test (with the accumulated server log) if it never comes up.
func waitForDemoHealth(t *testing.T, baseURL string) {
	t.Helper()

	healthURL := baseURL + "/health"
	deadline := time.Now().Add(demoStartTimeout)

	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(healthURL) //nolint:noctx,bodyclose // test fixture: URL is a local fixture server; body is tiny and closed below
		if err == nil {
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return
			}
		}

		time.Sleep(demoHealthPollInterval)
	}

	t.Fatalf("visualtest[demo]: server never became healthy at %s within %s", healthURL, demoStartTimeout)
}

// FailIfServerErrors fails the test when the demo server log contains an
// HTTP 500 response line. Used by the e2e suite so a handler crash inside an
// otherwise-passing browser flow cannot hide.
func (s *DemoServer) FailIfServerErrors(t *testing.T) {
	t.Helper()

	if bytes.Contains(s.Log.Bytes(), []byte("500 ")) {
		t.Fatalf("visualtest[demo]: server log contains 500 responses:\n%s", s.Log.String())
	}
}
