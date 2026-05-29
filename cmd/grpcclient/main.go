// Package main: minimal gRPC client that calls our server once and prints.
//
// In a real app the client lives INSIDE another service (e.g. your REST
// gateway calls the analytics gRPC server). This standalone client exists
// only as a "press the button, see it work" tool for the lesson.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/moonyango/taskbridge/proto/taskpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Step 1 — open a CONNECTION to the gRPC server.
	//
	// grpc.NewClient returns a ClientConn (a long-lived, multiplexed HTTP/2
	// connection). You typically create ONE per remote service and reuse it
	// for every call — it's safe for concurrent use.
	//
	// insecure.NewCredentials() means "no TLS" — fine for localhost. In prod
	// you'd use real TLS credentials.
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// Step 2 — wrap the connection in the TYPED CLIENT generated from the
	// .proto. From now on you call methods like a normal Go function — the
	// stub handles serialization, transport, deserialization for you.
	client := taskpb.NewTaskAnalyticsServiceClient(conn)

	// Step 3 — set a DEADLINE on the call. ALWAYS pass a context with a
	// timeout to RPC calls; without one a stuck server hangs your client
	// forever. 2 seconds is generous for a local call; in real services you
	// often tune this per-endpoint.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Step 4 — actually invoke the RPC. Looks like a local function call.
	// Notice we pass an empty taskpb.Empty — gRPC methods always require a
	// message argument, even when "no input" is the semantic.
	stats, err := client.GetTaskStats(ctx, &taskpb.Empty{})
	if err != nil {
		log.Fatalf("GetTaskStats: %v", err)
	}

	fmt.Printf("total=%d  completed=%d  pending=%d\n",
		stats.GetTotalTasks(),
		stats.GetCompletedTasks(),
		stats.GetPendingTasks(),
	)

	// Bonus call with a real argument.
	comp, err := client.GetTaskCompletion(ctx, &taskpb.TaskID{Id: 42})
	if err != nil {
		log.Fatalf("GetTaskCompletion: %v", err)
	}
	fmt.Printf("task id=%d  title=%q  completed=%v\n",
		comp.GetId(), comp.GetTitle(), comp.GetCompleted(),
	)
}
