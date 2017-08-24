package router

import (
	//"fmt"
	//"sync"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	//	"log"
	"fmt"
	"github.com/anemos-io/engine"
	"github.com/anemos-io/engine/provider/noop"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type observerServer struct {
	router *Router
}

type executorServer struct {
	instances chan api.TaskInstance
	router    *Router
}

type remoteObserverBinding struct {
	stream api.Observer_CommandStreamServer
}

type remoteExecutorBinding struct {
	stream api.Executor_CommandStreamServer
}

func (reb *remoteExecutorBinding) Execute(instance *api.TaskInstance) {
	command := api.ExecutorCommand{
		Instance: instance,
	}
	reb.stream.Send(&command)
}

func (s *observerServer) CommandStream(request *api.ObserverCommandStreamRequest, stream api.Observer_CommandStreamServer) error {
	executor := &remoteObserverBinding{
		stream: stream,
	}
	s.router.observers["anemos:noop/external"] = executor

	loop := true
	for loop {
		timeout := make(chan bool, 5)
		go func() {
			time.Sleep(1 * time.Second)
			timeout <- true
		}()

		select {
		//case cmd := <-s.instances:
		//	command := api.ExecutorCommand{
		//		Instance: &cmd,
		//	}
		//	stream.Send(&command)
		case <-timeout:
			log.Print("Timeout, nothing todo")
			// the read from ch has timed out
		}
		err := stream.Context().Err()
		if err != nil {
			log.Println(err.Error())
			loop = false
		}
	}

	return nil
}

func (s *observerServer) Trigger(ctx context.Context, in *api.TriggerRequest) (*api.TriggerResponse, error) {

	s.router.Trigger(in.Event)

	return nil, nil
}

func (s *executorServer) CommandStream(request *api.ExecutorCommandStreamRequest, stream api.Executor_CommandStreamServer) error {
	executor := &remoteExecutorBinding{
		stream: stream,
	}
	s.router.executors["anemos:noop/external"] = executor

	loop := true
	for loop {
		timeout := make(chan bool, 5)
		go func() {
			time.Sleep(1 * time.Second)
			timeout <- true
		}()

		select {
		//case cmd := <-s.instances:
		//	command := api.ExecutorCommand{
		//		Instance: &cmd,
		//	}
		//	stream.Send(&command)
		case <-timeout:
			log.Print("Timeout, nothing todo")
			// the read from ch has timed out
		}
		err := stream.Context().Err()
		if err != nil {
			log.Println(err.Error())
			loop = false
		}
	}

	return nil
}

type Router struct {
	Channel   chan *api.Event
	instances map[string]anemos.Node
	mutex     sync.Mutex

	executors map[string]anemos.Executor
	observers map[string]anemos.Observer

	executorServer *executorServer
	observerServer *observerServer
}

func NewRouter(server *grpc.Server) *Router {
	channel := make(chan *api.Event)

	executor := noop.NoopExecutor{}
	observer := noop.NoopObserver{}
	executor.CoupleObserver(&observer)

	router := &Router{
		Channel: channel,

		instances: make(map[string]anemos.Node),

		executors: make(map[string]anemos.Executor),
		observers: make(map[string]anemos.Observer),
	}
	observer.Router = router

	router.executors["anemos:noop"] = &executor
	router.observers["anemos:noop"] = &observer

	if server != nil {
		observerApi := &observerServer{}
		observerApi.router = router
		router.observerServer = observerApi
		executorApi := &executorServer{
			instances: make(chan api.TaskInstance),
		}
		executorApi.router = router
		router.executorServer = executorApi

		api.RegisterObserverServer(server, observerApi)
		api.RegisterExecutorServer(server, executorApi)
	}

	return router
}

func (r *Router) StartTask(node anemos.Node, instance *api.TaskInstance) {
	iid := fmt.Sprintf("%s:%s:%s:%s", instance.Provider, instance.Operation, instance.Name, instance.Id)
	log.Printf("Router.StartTask: iid(%s) and execute\n", iid)
	r.mutex.Lock()
	r.instances[iid] = node
	r.mutex.Unlock()

	e, found := r.executors["anemos:noop/external"]
	if found {
		e.Execute(instance)
	} else {
		r.executors["anemos:noop"].Execute(instance)
	}

	//event := api.Event{
	//	Uri: "anemos/event:manual",
	//}
	//node.
}

func (r *Router) StartVirtual(node anemos.Node, instance *api.TaskInstance) {
	iid := fmt.Sprintf("%s:%s:%s:%s", instance.Provider, instance.Operation, instance.Name, instance.Id)
	log.Printf("Router.StartVirtual: start iid(%s)\n", iid)
	r.mutex.Lock()
	r.instances[iid] = node
	r.mutex.Unlock()

	event := api.Event{
		Uri: anemos.Uri{
			Kind:      "anemos/event",
			Provider:  instance.Provider,
			Operation: instance.Operation,
			Name:      instance.Name,
			Id:        instance.Id,
			Status:    "finished",
		}.String(),
	}
	r.Trigger(&event)
}

func (r *Router) Fail(node anemos.Node, instance *api.TaskInstance) {
	iid := fmt.Sprintf("%s:%s:%s:%s", instance.Provider, instance.Operation, instance.Name, instance.Id)
	log.Printf("Router.Fail: start iid(%s)\n", iid)
	r.mutex.Lock()
	r.instances[iid] = node
	r.mutex.Unlock()

	event := api.Event{
		Uri: anemos.Uri{
			Kind:      "anemos/event",
			Provider:  instance.Provider,
			Operation: instance.Operation,
			Name:      instance.Name,
			Id:        instance.Id,
			Status:    "fail",
		}.String(),
	}
	r.Trigger(&event)
}

func (r *Router) SignalDownstream(node anemos.Node) {
	event := api.Event{
		Uri: anemos.Uri{
			Kind:      "anemos/event",
			Provider:  "anemos",
			Operation: "parent",
			Name:      node.Name(),
			Id:        "0000000000000000",
			Status:    "finished",
		}.String(),
	}
	go node.OnEvent(&event)
}

func (r *Router) Trigger(event *api.Event) {

	uri, _ := anemos.ParseUri(event.Uri)
	iid := fmt.Sprintf("%s:%s:%s:%s", uri.Provider, uri.Operation, uri.Name, uri.Id)
	log.Printf("Router.Trigger: iid(%s) and finish\n", iid)

	r.mutex.Lock()
	node := r.instances[iid]
	r.mutex.Unlock()

	node.OnFinish(event)

	//event := api.Event{
	//	Uri: "anemos/event:manual",
	//}
	//node.
}

func (r *Router) RegisterSession(session anemos.Session) {
	session.SetRouter(r)
}
