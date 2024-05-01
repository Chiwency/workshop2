package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	pb "github.com/chiwency/workshop2/rpc_gen"
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
)

const electionPath = "/election"

type masterServer struct {
	pb.UnimplementedMapReduceServiceServer
	isLeader bool
	workers  []string // list of wrker nodes
}

// TaskRequest handle task request from client
func (m *masterServer) SubmitTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
	log.Printf("Received task request: %v", in.GetTaskId())

	// if I am the leader, send the task to a worker
	if len(m.workers) > 0 {
		workerResponse, err := sendTaskToWorker(m.workers[0], in.GetTaskId(), in.GetMessage())
		if err != nil {
			return nil, err
		}
		return &pb.TaskResponse{Message: workerResponse}, nil
	}
	return nil, fmt.Errorf("no workers available")
}

// sendTaskToWorker sends a task to a worker node
func sendTaskToWorker(workerAddress, taskId string, msg string) (string, error) {
	conn, err := grpc.Dial(workerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return "", err
	}
	defer conn.Close()
	c := pb.NewMapReduceServiceClient(conn)

	// revoke the rpc call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.PerformTask(ctx, &pb.TaskRequest{TaskId: taskId, Message: msg})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return "", err
	}
	return r.GetMessage(), nil
}

var addr string

func main() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	addr = listener.Addr().String()
	grpcServer := grpc.NewServer()

	ms := &masterServer{workers: []string{"worker:7070"}}
	pb.RegisterMapReduceServiceServer(grpcServer, ms)

	// connect to ZooKeeper
	conn, _, err := zk.Connect([]string{"zookeeper:2181"}, time.Second*10)
	if err != nil {
		log.Fatalf("failed to connect to ZooKeeper: %v", err)
	}
	defer conn.Close()

	// make sure election path exists
	ensurePathExists(conn, electionPath)

	// join the election
	myNode, err := conn.CreateProtectedEphemeralSequential(electionPath+"/node", nil, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatalf("failed to create election node: %v", err)
	}
	fmt.Printf("Created election node: %s\n", myNode)

	go electLeader(conn, myNode, ms)

	log.Printf("Master server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func ensurePathExists(conn *zk.Conn, path string) {
	exists, _, err := conn.Exists(path)
	if err != nil {
		log.Fatalf("check path exists failed: %v", err)
	}
	if !exists {
		_, err := conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil && !errors.Is(err, zk.ErrNodeExists) {
			log.Fatalf("create path failed: %v", err)
		}
	}
}

func electLeader(conn *zk.Conn, myNode string, ms *masterServer) {
	for {
		children, _, err := conn.Children(electionPath)
		if err != nil {
			log.Printf("Failed to list children: %v", err)
			continue
		}

		sort.Strings(children)
		leader := electionPath + "/" + children[0]

		if leader == myNode {
			log.Println("I am the leader")
			if !ms.isLeader {
				log.Println("I became the leader")
				resp, _ := ms.SubmitTask(context.Background(), &pb.TaskRequest{TaskId: "leader_election", Message: addr})
				log.Printf("We got form worker: %s", resp.GetMessage())
				ms.isLeader = true
			}
		} else {
			log.Println("I am the follower")
			ms.isLeader = false
			time.Sleep(time.Second * 5) // Re-check after some delay
		}
	}
}
