// Package main is the entry point for the TaskBridge REST API binary.
//
// CONVENTION: anything under cmd/<name>/ compiles into a binary called <name>.
// We could later add cmd/grpc/main.go for a separate gRPC binary, and `go build ./...`
// would produce both. Keeping main.go TINY is the goal — it should only wire
// dependencies together and start the server. All real logic lives under internal/.
package main

import (
	"log"
	"net/http"
)

func main() {
	// http.NewServeMux is the standard-library router. We'll replace it with
	// Chi in Lesson 2 — Chi gives us URL params (/tasks/:id) and middleware
	// out of the box. Stdlib mux can do it too but it's clunkier.
	mux := http.NewServeMux()

	// A /health endpoint is the simplest "is the server alive?" probe.
	// Load balancers, Kubernetes, and Docker all expect one. Returning JSON
	// (not just "ok") makes it parseable by monitoring tools.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Listening on :8080 is convention for dev HTTP services in Go.
	// In production you'd read this from an env var or config file.
	addr := ":8080"
	log.Printf("TaskBridge listening on %s", addr)

	// ListenAndServe BLOCKS until the server errors. log.Fatal prints and
	// exits with code 1 — good enough for a CLI app; in prod you'd want
	// graceful shutdown via signal.Notify (we can add that if time permits).
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
