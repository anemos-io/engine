package router

import (
	//"fmt"
	//"sync"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	//	"log"
	"fmt"
	"github.com/anemos-io/engine"
	"github.com/anemos-io/engine/provider/noop"
	"log"
	"sync"
)

type InternalResourceRouter struct {
}

type InternalEventRouter struct {
}

type InternalRouter struct {
	Channel  chan *api.Event
	executor noop.NoopExecutor
	observer noop.NoopObserver

	*InternalResourceRouter
	*InternalEventRouter
	instances map[string]anemos.Node
	mutex     sync.Mutex
}

func (r *InternalRouter) ObserverLoop() {
	for true {
		event := <-r.observer.EventChannel
		r.Trigger(event)
	}
}

func NewInternalRouter() *InternalRouter {
	channel := make(chan *api.Event)

	executor := noop.NoopExecutor{}
	observer := noop.NoopObserver{
		EventChannel: channel,
	}
	executor.CoupleObserver(&observer)

	router := &InternalRouter{
		Channel:  channel,
		executor: executor,
		observer: observer,

		instances: make(map[string]anemos.Node),

		InternalEventRouter:    &InternalEventRouter{},
		InternalResourceRouter: &InternalResourceRouter{},
	}
	go router.ObserverLoop()
	return router
}

func (r *InternalRouter) StartTask(node anemos.Node, instance *api.TaskInstance) {
	iid := fmt.Sprintf("%s:%s:%s:%s", instance.Provider, instance.Operation, instance.Name, instance.Id)
	log.Printf("InternalRouter.StartTask: iid(%s) and execute\n", iid)
	r.mutex.Lock()
	r.instances[iid] = node
	r.mutex.Unlock()

	r.executor.Execute(instance)

	//event := api.Event{
	//	Uri: "anemos/event:manual",
	//}
	//node.
}

func (r *InternalRouter) StartVirtual(node anemos.Node, instance *api.TaskInstance) {
	iid := fmt.Sprintf("%s:%s:%s:%s", instance.Provider, instance.Operation, instance.Name, instance.Id)
	log.Printf("InternalRouter.StartVirtual: start iid(%s)\n", iid)
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

func (r *InternalRouter) Fail(node anemos.Node, instance *api.TaskInstance) {
	iid := fmt.Sprintf("%s:%s:%s:%s", instance.Provider, instance.Operation, instance.Name, instance.Id)
	log.Printf("InternalRouter.Fail: start iid(%s)\n", iid)
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

func (r *InternalRouter) SignalDownstream(node anemos.Node) {
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

func (r *InternalRouter) Trigger(event *api.Event) {

	uri, _ := anemos.ParseUri(event.Uri)
	iid := fmt.Sprintf("%s:%s:%s:%s", uri.Provider, uri.Operation, uri.Name, uri.Id)
	log.Printf("InternalRouter.Trigger: iid(%s) and finish\n", iid)

	r.mutex.Lock()
	node := r.instances[iid]
	r.mutex.Unlock()

	node.OnFinish(event)

	//event := api.Event{
	//	Uri: "anemos/event:manual",
	//}
	//node.
}

func (r *InternalRouter) RegisterSession(session anemos.Session) {
	session.SetRouter(r)
}
