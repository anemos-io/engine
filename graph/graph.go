package graph

import (
	//"fmt"
	//"sync"
	"github.com/anemos-io/engine"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"log"
	//"fmt"
	//"reflect"
	"fmt"
	"math/rand"
	"sync"
)

func logOn(node anemos.Node, function string, name string, event *api.Event) {
	log.Printf("%s[%s].%s: %s", function, node.Name(), name, event.Uri)
}

type Node struct {
	name       string
	session    anemos.Session
	downstream map[string]anemos.Node
	upstream   map[string]anemos.Node
	status     anemos.NodeInstanceStatus
	mutex      sync.Mutex
}

func NewNode() *Node {
	return &Node{
		downstream: make(map[string]anemos.Node, 0),
		upstream:   make(map[string]anemos.Node, 0),
		status:     anemos.Unknown,
		//nodes: make([]anemos.Node,0),
	}
}

func LinkDown(from anemos.Node, to anemos.Node) {
	name := fmt.Sprintf("%s>%s", from.Name(), to.Name())
	LinkDownNamed(from, to, name)
}

func LinkDownNamed(from anemos.Node, to anemos.Node, name string) {
	from.AddDownstream(name, to)
	to.AddUpstream(name, from)
	log.Printf("Link: %s > %s", from.Name(), to.Name())
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) AddUpstream(name string, node anemos.Node) {
	n.upstream[name] = node
}

func (n *Node) AddDownstream(name string, node anemos.Node) {
	n.downstream[name] = node
}

func (n *Node) Downstream() map[string]anemos.Node {
	return n.downstream
}

func (n *Node) Upstream() map[string]anemos.Node {
	return n.upstream
}

func (n *Node) Status() anemos.NodeInstanceStatus {
	return n.status
}

func (n *Node) AssignSession(session anemos.Session) {
	n.session = session
}

func NewTaskNode() *TaskNode {
	return &TaskNode{
		Node:       NewNode(),
		attributes: make(map[string]string, 0),
	}
}

type TaskNode struct {
	*Node
	provider   string
	operation  string
	attributes map[string]string
}

func (n *TaskNode) Provider() string {
	return n.provider
}

func (n *TaskNode) Operation() string {
	return n.operation
}

func (n *TaskNode) Attributes() map[string]string {
	return n.attributes
}

func (n *TaskNode) OnEvent(event *api.Event) {
	logOn(n, "TaskNode", "OnEvent", event)
	status := anemos.Success
	if len(n.upstream) > 0 {
		for _, node := range n.upstream {
			if !node.EndStateReached() {
				log.Printf("TaskNode.OnEvent: Not all dependencies are saticfied for %s", n.name)
				return
			} else {
				if node.Status() == anemos.Fail && status < anemos.Fail {
					status = anemos.Fail
				} else if node.Status() == anemos.Skip && status < anemos.Skip {
					// TODO: Skip? Defaults to success
				}
			}
		}
	}

	if status == anemos.Success {
		def := n.session.NewTaskInstance(n)
		log.Printf("TaskNode[%s].OnEvent: All dependencies are saticfied.", n.name)

		n.session.Router().StartTask(n, def)
	} else {
		log.Printf("TaskNode[%s].OnEvent: Dependency have failure, fail and finish.", n.name)
		def := api.TaskInstance{
			Provider:  "anemos",
			Operation: "virtual",
			Name:      n.name,
			Id:        fmt.Sprintf("%x", rand.Int63()),
		}
		n.session.Router().Fail(n, &def)
	}
}

func (n *TaskNode) OnFinish(event *api.Event) {
	logOn(n, "TaskNode", "OnFinish", event)
	if n.status == anemos.Success {
		log.Printf("TaskNode.OnFinish: WARNING: Already succeeded")
		return
	}

	eventUri, _ := anemos.ParseUri(event.Uri)
	status := anemos.Success
	if eventUri.Status == "success" {
		n.status = anemos.Success
	} else if eventUri.Status == "finished" {
		// rules
		status = anemos.Success
		if len(n.upstream) > 0 {
			for _, node := range n.upstream {
				if node.Status() == anemos.Fail {
					status = anemos.Fail
				}
				//if !node.EndStateReached() {
				//	log.Printf("VirtualNode.OnEvent: Not all dependencies are saticfied for %s", n.name)
				//	return
				//}
			}
		}
	} else if eventUri.Status == "fail" {
		status = anemos.Fail
	}
	n.status = status

	if n.EndStateReached() {
		if len(n.downstream) > 0 {
			for _, node := range n.downstream {
				n.session.Router().SignalDownstream(node)
			}
		}
	}
}

func (n *TaskNode) OnStart(event *api.Event) {
	logOn(n, "TaskNode", "OnStart", event)

}

func (n *TaskNode) OnProgress(event *api.Event) {
	logOn(n, "TaskNode", "OnProgress", event)

}

func (n *TaskNode) OnCancel(event *api.Event) {
	logOn(n, "TaskNode", "OnCancel", event)

}

func (n *TaskNode) OnSkip(event *api.Event) {
	logOn(n, "TaskNode", "OnSkip", event)

}

func (n *TaskNode) EndStateReached() bool {
	return n.status == anemos.Success || n.status == anemos.Fail || n.status == anemos.Skip
}

func NewVirtualNode() *VirtualNode {
	return &VirtualNode{
		Node: NewNode(),
		//nodes: make([]anemos.Node,0),
	}
}

type VirtualNodeType int

const (
	Solo VirtualNodeType = iota
	Begin
	End
)

type VirtualNode struct {
	*Node
	parent *Group
	kind   VirtualNodeType
}

func (n *VirtualNode) Provider() string {
	return "anemos"
}

func (n *VirtualNode) Operation() string {
	return "virtual"
}

func (n *VirtualNode) Attributes() map[string]string {
	return nil
}

func (n *VirtualNode) OnEvent(event *api.Event) {
	logOn(n, "VirtualNode", "OnEvent", event)

	if len(n.upstream) > 0 {
		for _, node := range n.upstream {
			if !node.EndStateReached() {
				log.Printf("VirtualNode.OnEvent: Not all dependencies are saticfied for %s", n.name)
				return
			}
		}
	}

	def := api.TaskInstance{
		Provider:  "anemos",
		Operation: "virtual",
		Name:      n.name,
		Id:        fmt.Sprintf("%x", rand.Int63()),
	}

	log.Printf("VirtualNode.OnEvent: All dependencies are saticfied for %s", n.name)
	n.session.Router().StartVirtual(n, &def)
}

func (n *VirtualNode) OnFinish(event *api.Event) {
	logOn(n, "VirtualNode", "OnFinish", event)
	if n.status == anemos.Success {
		log.Printf("VirtualNode.OnFinish: WARNING Already succeeded")
		return
	}

	// rules
	status := anemos.Success
	if len(n.upstream) > 0 {
		for _, node := range n.upstream {
			if node.Status() == anemos.Fail {
				status = anemos.Fail
			}
			//if !node.EndStateReached() {
			//	log.Printf("VirtualNode.OnEvent: Not all dependencies are saticfied for %s", n.name)
			//	return
			//}
		}
	}

	// eventUri, _ := anemos.ParseUri(event.Uri)
	// status

	n.status = status
	if len(n.downstream) > 0 {
		for _, node := range n.downstream {
			n.session.Router().SignalDownstream(node)
		}
	}
	if n.parent != nil && n.kind == End {
		log.Printf("VirtualNode[%s].OnFinish: End node reached for group %s", n.Name(),
			n.parent.Name())
		n.parent.status = n.status
		n.parent.channel <- n.status == anemos.Success
	}
}

func (n *VirtualNode) OnStart(event *api.Event) {
	logOn(n, "VirtualNode", "OnStart", event)

}

func (n *VirtualNode) OnProgress(event *api.Event) {
	logOn(n, "VirtualNode", "OnProgress", event)

}

func (n *VirtualNode) OnCancel(event *api.Event) {
	logOn(n, "VirtualNode", "OnCancel", event)

}

func (n *VirtualNode) OnSkip(event *api.Event) {
	logOn(n, "VirtualNode", "OnSkip", event)

}

func (n *VirtualNode) EndStateReached() bool {
	return n.status == anemos.Success || n.status == anemos.Fail
}

type Group struct {
	*Node
	nodes   []anemos.Node
	begin   *VirtualNode
	end     *VirtualNode
	channel chan bool
}

func NewGroup() *Group {
	return &Group{
		Node:    &Node{},
		nodes:   make([]anemos.Node, 0),
		channel: make(chan bool),
	}
}

func (n *Group) Provider() string {
	return "anemos"
}

func (n *Group) Operation() string {
	return "group"
}

func (n *Group) Attributes() map[string]string {
	return nil
}

func CopyGroup(source *Group) *Group {

	return nil
}

func (g *Group) Resolve() {
	g.begin = NewVirtualNode()
	g.begin.parent = g
	g.begin.kind = Begin
	g.begin.name = g.name + "+begin"
	g.end = NewVirtualNode()
	g.end.parent = g
	g.end.kind = End
	g.end.name = g.name + "+end"

	for _, node := range g.nodes {
		if len(node.Upstream()) == 0 {
			name := fmt.Sprintf("%s>%s", g.begin.Name(), node.Name())
			log.Printf("Group Resolver: Adding link for %s", name)
			LinkDown(g.begin, node)
		}
		if len(node.Downstream()) == 0 {
			name := fmt.Sprintf("%s>%s", node.Name(), g.end.Name())
			log.Printf("Group Resolver: Adding link for %s", name)
			LinkDown(node, g.end)
		}

	}
}

func (g *Group) AddNode(node anemos.Node) {
	g.nodes = append(g.nodes, node)
}

func (g *Group) AssignSession(session anemos.Session) {
	g.session = session
	g.begin.session = session
	g.end.session = session
	for _, node := range g.nodes {
		node.AssignSession(session)
	}
}

func (g *Group) OnEvent(event *api.Event) {
	logOn(g, "Group", "OnEvent", event)
	g.begin.OnEvent(event)
}

func (g *Group) OnStart(event *api.Event) {
	logOn(g, "Group", "OnStart", event)

}

func (g *Group) OnProgress(event *api.Event) {
	logOn(g, "Group", "OnProgress", event)

}

func (g *Group) OnFinish(event *api.Event) {
	logOn(g, "Group", "OnFinish", event)

}

func (g *Group) OnCancel(event *api.Event) {
	logOn(g, "Group", "OnCancel", event)

}

func (g *Group) OnSkip(event *api.Event) {
	logOn(g, "Group", "OnSkip", event)

}

func (n *Group) EndStateReached() bool {
	return false
}
