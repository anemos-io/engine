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
	"fmt"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"github.com/anemos-io/engine/provider/noop"
	"io"
	"log"
)

const (
	port = ":5000"
)

type Server struct {
}

type observerBinding struct {
	observerClient api.ObserverClient
}

func (ob *observerBinding) Trigger(e *api.Event) {
	request := &api.TriggerRequest{
		Event: e,
	}

	ob.observerClient.Trigger(context.Background(), request)
}

func main() {

	conn, err := grpc.Dial("localhost:4270", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer conn.Close()

	executor := noop.NoopExecutor{}
	observer := noop.NoopObserver{}

	binding := &observerBinding{}
	binding.observerClient = api.NewObserverClient(conn)
	observer.Router = binding

	executor.CoupleObserver(&observer)
	executorClient := api.NewExecutorClient(conn)

	request := api.ExecutorCommandStreamRequest{}

	stream, err := executorClient.CommandStream(context.Background(), &request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for {
		command, err := stream.Recv()
		fmt.Println("Received")
		fmt.Println(request)
		executor.Execute(command.Instance)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.CommandStream(_) = _, %v", executor, err)
		}
		log.Println(command)
	}

}
