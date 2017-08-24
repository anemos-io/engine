package main

import (
	"log"
	//"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	//"github.com/golang/protobuf/proto"
	//wkt_timestamp "github.com/golang/protobuf/ptypes/timestamp"
	//wkt_empty "github.com/golang/protobuf/ptypes/empty"
	"fmt"
	"github.com/anemos-io/engine"
	"github.com/anemos-io/engine/graph"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"github.com/anemos-io/engine/router"
	"github.com/anemos-io/engine/store"
	"time"
)

const (
	port = ":4270"
)

func test(st *store.GraphStore, router anemos.Router) {
	tc := time.Tick(10 * time.Second)

	for {
		<-tc
		group := st.Graphs["three-task-simple"]
		session := graph.NewSession(group)
		group.AssignSession(session)
		router.RegisterSession(session)

		event := api.Event{
			Uri: "anemos/event:manual",
		}
		g := session.Graph

		go g.OnEvent(&event)
	}
}

func main() {

	st := store.NewGraphStore()
	fmt.Println(st)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	router := router.NewRouter(s)

	go test(st, router)

	// Register reflection service on gRPC.
	reflection.Register(s)
	log.Println("Starting Anemos Engine")

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
