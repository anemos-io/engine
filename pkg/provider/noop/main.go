package main

import (
	//"log"
	"golang.org/x/net/context"
	//"net"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	//"github.com/golang/protobuf/proto"
	//wkt_timestamp "github.com/golang/protobuf/ptypes/timestamp"
	//wkt_empty "github.com/golang/protobuf/ptypes/empty"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"log"
	"io"
	"fmt"
)

const (
	port = ":5000"
)

type Server struct {
}

//func (s *Server) PostEvents(ctx context.Context, in *pb.PostEventsRequest) (*wkt_empty.Empty, error) {
//
//}

func main() {

	conn, err := grpc.Dial("localhost:4270", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer conn.Close()
	executor := api.NewExecutorClient(conn)

	request := api.ExecutorCommandStreamRequest{
	}

	stream, err := executor.CommandStream(context.Background(),&request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for {
		command, err := stream.Recv()
		fmt.Println("Received")
		fmt.Println(request)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.CommandStream(_) = _, %v", executor, err)
		}
		log.Println(command)
	}


}
