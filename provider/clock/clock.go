package clock

import (
	"container/list"
	"time"
	"fmt"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
)

type ClockObserver struct {
	channel chan *api.Event
	ticks   list.List
	looping bool
}

type ClockTaskDefinition struct {
	tick time.Time
	uri  string
}

func (c *ClockObserver) Tick() bool {

	for e := c.ticks.Front(); e != nil; e = e.Next() {
		val := e.Value.(ClockTaskDefinition)
		if time.Now().After(val.tick) {
			c.Trigger(e)
			fmt.Println("OK")
		}
	}
	return true
}

func (c *ClockObserver) Loop() bool {
	for ; c.looping; {
		c.Tick()
		time.Sleep(1000 * time.Millisecond)
	}

	return true
}

func (c *ClockObserver) Start() {
	c.looping = true
	go c.Loop()
}

func (c *ClockObserver) Stop() {
	c.looping = false
}

func (c *ClockObserver) Add(uri string, t time.Time) {
	event := ClockTaskDefinition{
		tick: t,
		uri:  uri,
	}
	c.ticks.PushBack(event)
}

func (c *ClockObserver) Trigger(element *list.Element) {

	val := element.Value.(ClockTaskDefinition)

	event := api.Event{
		Uri:"anemos/event:clock:0000:trigger",
		Metadata:make(map[string]string),
	}
	event.Metadata["anemos/metadata:task:timestamp"] = time.Now().Format(time.RFC3339Nano)
	c.channel <- &event

	fmt.Println(val.uri)
	c.ticks.Remove(element)
}
