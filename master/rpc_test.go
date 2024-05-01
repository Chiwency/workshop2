package main

import (
	"context"
	pb "github.com/chiwency/workshop2/rpc_gen"
	"google.golang.org/grpc"
	"log"
	"testing"
)

// mock a master server which sneds a message to worker
func TestPerformTask(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:7070", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewMapReduceServiceClient(conn)

	taskID := "1"
	resp, err := client.PerformTask(ctx, &pb.TaskRequest{TaskId: taskID, Message: "Hello, worker!"})
	if err != nil {
		t.Errorf("PerformTask failed: %v", err)
	}
	log.Printf("We got: %s", resp.GetMessage())
}
