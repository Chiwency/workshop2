package main

import (
	"context"
	pb "github.com/chiwency/workshop2/rpc_gen"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedMapReduceServiceServer
}

func (s *server) SubmitTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
	log.Printf("Received: %v", in.GetTaskId())
	// code for task submission
	return &pb.TaskResponse{Message: "Completed task " + in.GetTaskId()}, nil
}

func (s *server) PerformTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
	log.Printf("Performing task: %v", in.GetTaskId())
	// code for task execution
	return &pb.TaskResponse{Message: in.GetMessage() + " gatech"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":7070")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMapReduceServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
