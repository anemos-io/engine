package anemos

import (
	//"fmt"
	//"sync"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	//"log"
	"strings"
)

type NodeInstanceStatus int

const (
	Unknown     NodeInstanceStatus = iota
	Retry
	Initialized
	Queue
	Start
	Success
	Skip
	Fail
)

type Node interface {
	ShortName() (string)
	//Name() (string)
	AddUpstream(name string, node Node) ()
	AddDownstream(name string, node Node) ()
	Upstream() (map[string]Node)
	Downstream() (map[string]Node)
	Status() (NodeInstanceStatus)
	EndStateReached() (bool)
	SetRouter(router Router) ()

	OnEvent(event *api.Event)
	OnStart(event *api.Event)
	OnProgress(event *api.Event)
	OnFinish(event *api.Event)
	OnCancel(event *api.Event)
	OnSkip(event *api.Event)
}

type Group interface {
	Node
	AddNode(node Node) ()
	Resolve()
}

type ResourceRouter interface {
	Start(node Node, instance *api.TaskInstance)
	StartVirtual(node Node, instance *api.TaskInstance)
	Fail(node Node, instance *api.TaskInstance)
}

type EventRouter interface {
}

type Router interface {
	ResourceRouter
	EventRouter
	RegisterGroup(group Group)
}

type Uri struct {
	Kind      string
	Provider  string
	Operation string
	Name      string
	Id        string
	Status    string
}

func ParseUri(us string) (*Uri, error) {
	parts := strings.Split(us, ":")
	return &Uri{
		Kind:      parts[0],
		Provider:  parts[1],
		Operation: parts[2],
		Name:      parts[3],
		Id:        parts[4],
		Status:    parts[5],
	}, nil
}

const (
	MetaTaskRetry      = "anemos/meta:anemos:task:retry"
	MetaEventTimestamp = "anemos/meta:anemos:event:timestamp"
)
