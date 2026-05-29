// Package grpcserver implements the gRPC TaskAnalyticsService.
//
// What lives here vs in cmd/grpcsrv/main.go:
//   - HERE: the BUSINESS LOGIC of each RPC (how stats get computed).
//   - cmd/grpcsrv/main.go: the WIRING (sockets, registration, Serve).
// Keep these separate so you can unit-test the server WITHOUT opening a port.
package grpcserver

import (
	"context"

	"github.com/moonyango/taskbridge/proto/taskpb"
)

// Server is OUR implementation of the generated TaskAnalyticsServiceServer
// interface. Naming convention: most Go projects just call it "Server".
//
// THE EMBED IS NOT OPTIONAL:
// We embed taskpb.UnimplementedTaskAnalyticsServiceServer (by value, not
// pointer — see the comment in the generated file). This gives us:
//   1. Default "Unimplemented" stubs for any RPC we haven't written yet,
//      so the code compiles even mid-development.
//   2. FORWARD COMPATIBILITY — if someone later adds a new rpc to the
//      .proto, our Server still compiles; calls to the new rpc just return
//      codes.Unimplemented until we override the method here.
type Server struct {
	taskpb.UnimplementedTaskAnalyticsServiceServer

	// In a real app, we'd hold a reference to the repository or service
	// layer here, e.g. `repo TaskRepository`. The gRPC server would call
	// into the SAME service layer the REST handlers use — that's the
	// whole reason we keep handlers thin: one service, multiple transports.
}

// NewServer is the constructor. Trivial today, but in real code it'd take
// dependencies (a repo, a logger, metrics, etc.) so callers can swap them in
// tests. NEVER reach for package-level globals — pass deps through here.
func NewServer() *Server {
	return &Server{}
}

// GetTaskStats is one of our two RPC methods.
//
// SIGNATURE NOTES — these match exactly what the generated SERVER interface
// requires (no `opts...`, just ctx + request):
//
//	GetTaskStats(context.Context, *taskpb.Empty) (*taskpb.TaskStats, error)
//
// The ctx carries:
//   - request DEADLINE  (so we can give up if the client gave up)
//   - request METADATA  (headers — auth tokens, trace IDs)
//   - CANCELLATION      (caller closed the connection)
// Always propagate ctx into any downstream call (DB, HTTP, other gRPC).
func (s *Server) GetTaskStats(ctx context.Context, _ *taskpb.Empty) (*taskpb.TaskStats, error) {
	// Hardcoded for now — we don't have a database. In Phase 4 of the real
	// spec, these numbers would come from the repository/service layer:
	//   stats, err := s.svc.ComputeStats(ctx)
	//   if err != nil { return nil, status.Error(codes.Internal, err.Error()) }
	//   return &taskpb.TaskStats{ TotalTasks: stats.Total, ... }, nil
	return &taskpb.TaskStats{
		TotalTasks:     10,
		CompletedTasks: 7,
		PendingTasks:   3,
	}, nil
}

// GetTaskCompletion shows the pattern when there IS a request body.
// Notice we use the generated nil-safe getter (req.GetId()) rather than
// req.Id directly — habit worth forming, prevents nil-pointer panics if a
// future caller forgets to send the request.
func (s *Server) GetTaskCompletion(ctx context.Context, req *taskpb.TaskID) (*taskpb.TaskCompletion, error) {
	return &taskpb.TaskCompletion{
		Id:        req.GetId(),
		Completed: false,
		Title:     "Finish gRPC lesson",
	}, nil
}
