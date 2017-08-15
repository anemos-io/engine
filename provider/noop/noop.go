package noop

import (
	"container/list"
	"time"
	//"fmt"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"github.com/anemos-io/engine"
	"strconv"
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

const (
	AttrSuccessDuration = "anemos/attribute:anemos:noop:duration/success"
	AttrFailDuration    = "anemos/attribute:anemos:noop:duration/fail"
	AttrRetries         = "anemos/attribute:anemos:noop:reties"
	AttrCouple          = "anemos/attribute:anemos:noop:couple"
)

type MetaData string

type NoopTask struct {
	noopTaskType NoopTaskType
	duration     time.Duration
}

type NoopTaskDefinition struct {
	instance *api.TaskInstance

	retries          int
	success_duration time.Duration
	fail_duration    time.Duration
	couple           bool

	retry int
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

func (ne *NoopExecutor) execute(definition NoopTaskDefinition) {
	time.Sleep(definition.success_duration)
	ne.observer.Trigger(&definition)
}

func (ne *NoopExecutor) Execute(instance *api.TaskInstance) {

	definition := NoopTaskDefinition{
		instance: instance,
		//task: NoopTask{
		//	noopTaskType: Success,
		//	duration:     time.Duration(1 * time.Millisecond),
		//},
	}

	value, _ := instance.Attributes[AttrRetries]
	definition.retries, _ = strconv.Atoi(value)
	value, _ = instance.Attributes[AttrSuccessDuration]
	definition.success_duration, _ = time.ParseDuration(value)
	value, _ = instance.Attributes[AttrFailDuration]
	definition.fail_duration, _ = time.ParseDuration(value)

	value, _ = instance.Attributes[anemos.MetaRetry]
	definition.retry, _ = strconv.Atoi(value)

	go ne.execute(definition)

}
