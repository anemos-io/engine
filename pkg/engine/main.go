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
	"github.com/anemos-io/engine/router"
)

const (
	port = ":4270"
)

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	router.NewRouter(s)

	// Register reflection service on gRPC.
	reflection.Register(s)
	log.Println("Starting Anemos Engine")




	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
