package main

import (
	"log"
	//"golang.org/x/net/context"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	//"github.com/golang/protobuf/proto"
	//wkt_timestamp "github.com/golang/protobuf/ptypes/timestamp"
	//wkt_empty "github.com/golang/protobuf/ptypes/empty"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"time"
)

const (
	port = ":5000"
)

type observerServer struct {
}

func (s *observerServer) CommandStream(request *api.ObserverCommandStreamRequest, stream api.Observer_CommandStreamServer) error {
	return nil
}

type executorServer struct {
}

func (s *executorServer) CommandStream(request *api.ExecutorCommandStreamRequest, stream api.Executor_CommandStreamServer) error {
	for {
		command := api.ExecutorCommand{

		}

		stream.Send(&command)
		time.Sleep(1 * time.Second)
	}



	return nil
}

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterObserverServer(s, &observerServer{})
	api.RegisterExecutorServer(s, &executorServer{})
	// Register reflection service on gRPC.
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
