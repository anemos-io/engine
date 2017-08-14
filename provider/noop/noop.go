package noop

import (
	"container/list"
	"time"
	//"fmt"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"github.com/anemos-io/engine"
)

type NoopObserver struct {
	EventChannel chan *api.Event
	ticks        list.List
	looping      bool
}

type NoopTaskType int

const (
	Success NoopTaskType = iota
	Fail
)

type NoopTask struct {
	noopTaskType NoopTaskType
	duration     time.Duration
}

type NoopTaskDefinition struct {
	instance *api.TaskInstance
	task     NoopTask
	retry    []NoopTask
}

type NoopExecutor struct {
	observer *NoopObserver
}

func (c *NoopObserver) Start() {
	c.looping = true
	//go c.Loop()
}

func (c *NoopObserver) Stop() {
	c.looping = false
}

func (c *NoopObserver) Trigger(definition *NoopTaskDefinition) {

	event := api.Event{
		Uri: anemos.Uri{
			Kind:      "anemos/event",
			Provider:  "anemos",
			Operation: "noop",
			Name:      definition.instance.Name,
			Id:        definition.instance.Id,
			Status:    "success",
		}.String(),
		Metadata: make(map[string]string),
	}
	event.Metadata["anemos/metadata:task:timestamp"] = time.Now().Format(time.RFC3339Nano)
	c.EventChannel <- &event
}

func (ne *NoopExecutor) CoupleObserver(observer *NoopObserver) {
	ne.observer = observer
}

func (ne *NoopExecutor) _Execute(definition NoopTaskDefinition) {
	time.Sleep(definition.task.duration)
	ne.observer.Trigger(&definition)
}

func (ne *NoopExecutor) ExecuteTMP(definition NoopTaskDefinition) {

	go ne._Execute(definition)

}

func (ne *NoopExecutor) Execute(instance *api.TaskInstance) {

	definition := NoopTaskDefinition{
		instance: instance,
		task: NoopTask{
			noopTaskType: Success,
			duration:     time.Duration(1 * time.Millisecond),
		},
	}

	go ne._Execute(definition)

}
