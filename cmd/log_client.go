package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "oversee/agent"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:4092", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAgentClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.Log(ctx, &pb.LogRequest{
		ServiceName:       "Auditable",
		Operation:         "demo_log_audit",
		ActorId:           "demo_app",
		ActorType:         "user",
		AffectedResources: []string{"logs"},
		Metadata:          nil,
		IntegrityHash:     "",
		Timestamp:         timestamppb.Now(),
	})

	if err != nil {
		log.Fatalf("could not log: %v", err)
	}
}
