package router

import (
	//"fmt"
	//"sync"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	//	"log"
	"fmt"
	"github.com/anemos-io/engine"
	"github.com/anemos-io/engine/provider/noop"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type observerServer struct {
	router *Router
}

func (s *observerServer) CommandStream(request *api.ObserverCommandStreamRequest, stream api.Observer_CommandStreamServer) error {
	return nil
}

type executorServer struct {
	instances chan api.TaskInstance

	router *Router
}

func (s *executorServer) CommandStream(request *api.ExecutorCommandStreamRequest, stream api.Executor_CommandStreamServer) error {
	for {
		timeout := make(chan bool, 5)
		go func() {
			time.Sleep(1 * time.Second)
			timeout <- true
		}()

		select {
		case cmd := <-s.instances:
			command := api.ExecutorCommand{
				Instance: &cmd,
			}
			stream.Send(&command)
		case <-timeout:
			log.Print("Timeout, nothing todo")
			// the read from ch has timed out
		}
	}

	return nil
}

type Router struct {
	Channel  chan *api.Event
	executor noop.NoopExecutor
	observer noop.NoopObserver

	instances map[string]anemos.Node
	mutex     sync.Mutex

	executorServer *executorServer
	observerServer *observerServer
}

func (r *Router) ObserverLoop() {
	for true {
		event := <-r.observer.EventChannel
		r.Trigger(event)
	}
}

func NewRouter(server *grpc.Server) *Router {
	channel := make(chan *api.Event)

	executor := noop.NoopExecutor{}
	observer := noop.NoopObserver{
		EventChannel: channel,
	}
	executor.CoupleObserver(&observer)

	router := &Router{
		Channel:  channel,
		executor: executor,
		observer: observer,

		instances: make(map[string]anemos.Node),
	}

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

	go router.ObserverLoop()
	return router
}

func (r *Router) StartTask(node anemos.Node, instance *api.TaskInstance) {
	iid := fmt.Sprintf("%s:%s:%s:%s", instance.Provider, instance.Operation, instance.Name, instance.Id)
	log.Printf("Router.StartTask: iid(%s) and execute\n", iid)
	r.mutex.Lock()
	r.instances[iid] = node
	r.mutex.Unlock()

	r.executor.Execute(instance)

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
