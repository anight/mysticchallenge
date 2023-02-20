package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	p "github.com/anight/mysticchallenge/proto"
	"google.golang.org/grpc"
)

var (
	grpcPort = flag.Int("grpcport", 8000, "The grpc server port")
	workers  = flag.Int("workers", 10, "The number of workers")
)

type server struct {
	p.UnimplementedRemoteExecuteAPIServer
	j *JobServer
}

func (s *server) Execute(ctx context.Context, req *p.RequestExecute) (*p.ResponseExecute, error) {
	result, err := s.j.Execute(req.Request)
	if err != nil {
		return &p.ResponseExecute{Error: fmt.Sprintf("Execute() failed: %v", err)}, nil
	}
	return &p.ResponseExecute{Result: result}, nil
}

func (s *server) GetWorkers(ctx context.Context, req *p.RequestGetWorkers) (*p.ResponseGetWorkers, error) {
	return &p.ResponseGetWorkers{Workers: int32(*workers)}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	s := &server{j: NewJobServer(*workers)}
	p.RegisterRemoteExecuteAPIServer(grpcServer, s)
	log.Printf("GRPC: server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("GRPC: failed to serve: %v", err)
	}
}
