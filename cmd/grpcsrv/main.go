// Package main: the gRPC server binary.
//
// CONVENTION: anything under cmd/<name>/main.go produces a binary called <name>.
// This file's job is small: open a TCP socket, create a gRPC server, plug in
// our implementation, and call Serve. All real logic lives in internal/grpcserver.
package main

import (
	"log"
	"net"

	"github.com/moonyango/taskbridge/internal/grpcserver"
	"github.com/moonyango/taskbridge/proto/taskpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Step 1 — open a raw TCP socket.
	//
	// :50051 is the CONVENTIONAL default port for gRPC services (the way :80
	// is for HTTP). You can use anything; this is just what most tutorials,
	// Docker images, and tooling defaults expect.
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}

	// Step 2 — create the gRPC SERVER object.
	//
	// grpc.NewServer() returns a server with no services registered yet.
	// You can pass options here (interceptors, max message size, TLS creds);
	// for now we want the bare default.
	s := grpc.NewServer()

	// Step 3 — register our IMPLEMENTATION against the SERVICE DESCRIPTOR.
	//
	// `RegisterTaskAnalyticsServiceServer` is a GENERATED function (lives in
	// task_grpc.pb.go). It tells the gRPC runtime: "when you receive a call
	// for service `task.TaskAnalyticsService`, dispatch it to this Server
	// instance." Without this call, the gRPC server would receive bytes and
	// have nowhere to route them.
	taskpb.RegisterTaskAnalyticsServiceServer(s, grpcserver.NewServer())

	// Step 4 (optional but helpful) — enable gRPC reflection.
	//
	// Reflection lets tools like `grpcurl` discover available methods at
	// runtime without needing the .proto file. Like swagger-UI for gRPC.
	// Don't enable this on a production internet-facing server; for internal
	// services it's a quality-of-life win.
	reflection.Register(s)

	log.Printf("gRPC TaskAnalyticsService listening on %s", addr)

	// Step 5 — Serve BLOCKS forever, accepting connections and dispatching.
	// In production you'd add a signal handler and call s.GracefulStop()
	// on SIGTERM. We'll skip that for the lesson.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
