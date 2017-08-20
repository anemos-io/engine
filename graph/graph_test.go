package graph

import (
	"github.com/anemos-io/engine"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"github.com/anemos-io/engine/provider/noop"
	"github.com/anemos-io/engine/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewSuccessTask(name string) *TaskNode {
	node := NewTaskNode()
	node.provider = "anemos"
	node.operation = "noop"
	node.name = name

	return node
}

func StartSessionForSuccess(session *Session, expectSuccess bool, t *testing.T) {
	event := api.Event{
		Uri: "anemos/event:manual",
	}
	g := session.graph

	go g.OnEvent(&event)
	assert.Equal(t, expectSuccess, <-g.channel)
}

func TestTwoTasks(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task2 := NewSuccessTask("task2")
	LinkDown(task1, task2)

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.AddNode(task2)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, true, t)

	assert.Equal(t, anemos.Success, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Success, s.graph.Status())
}

func TestSimpleSplit(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task2 := NewSuccessTask("task2")
	task3 := NewSuccessTask("task3")
	LinkDown(task1, task2)
	LinkDown(task1, task3)

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.AddNode(task2)
	g.AddNode(task3)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, true, t)

	assert.Equal(t, anemos.Success, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task3"].Status())
	assert.Equal(t, anemos.Success, s.graph.Status())
}

func TestSimpleJoin(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task2 := NewSuccessTask("task2")
	task3 := NewSuccessTask("task3")
	LinkDown(task1, task3)
	LinkDown(task2, task3)

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.AddNode(task2)
	g.AddNode(task3)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, true, t)

	assert.Equal(t, anemos.Success, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task3"].Status())
	assert.Equal(t, anemos.Success, s.graph.Status())
}

func TestSimpleSplitAndJoin(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task2 := NewSuccessTask("task2")
	task3 := NewSuccessTask("task3")
	task4 := NewSuccessTask("task4")
	LinkDown(task1, task2)
	LinkDown(task1, task3)
	LinkDown(task2, task4)
	LinkDown(task3, task4)

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.AddNode(task2)
	g.AddNode(task3)
	g.AddNode(task4)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, true, t)

	assert.Equal(t, anemos.Success, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task3"].Status())
	assert.Equal(t, anemos.Success, s.graph.nodes["task4"].Status())
	assert.Equal(t, anemos.Success, s.graph.Status())
}

func TestSingleFail(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task1.attributes[noop.AttrNameRetries] = "1"

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, false, t)

	assert.Equal(t, anemos.Fail, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Fail, s.graph.Status())
}

func TestFailPropagation(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task1.attributes[noop.AttrNameRetries] = "1"
	task2 := NewSuccessTask("task2")
	//task2.attributes[noop.AttrNameRetries] = "1"
	LinkDown(task1, task2)

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.AddNode(task2)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, false, t)

	assert.Equal(t, anemos.Fail, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Fail, s.graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Fail, s.graph.Status())
}

func TestFailSplitPropagation(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task1.attributes[noop.AttrNameRetries] = "1"
	task2 := NewSuccessTask("task2")
	task3 := NewSuccessTask("task3")
	//task2.attributes[noop.AttrNameRetries] = "1"
	LinkDown(task1, task2)
	LinkDown(task1, task3)

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.AddNode(task2)
	g.AddNode(task3)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, false, t)

	assert.Equal(t, anemos.Fail, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Fail, s.graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Fail, s.graph.nodes["task3"].Status())
	assert.Equal(t, anemos.Fail, s.graph.Status())
}

func TestFailSplitJoinPropagation(t *testing.T) {

	r := router.NewInternalRouter()

	task1 := NewSuccessTask("task1")
	task1.attributes[noop.AttrNameRetries] = "1"
	task2 := NewSuccessTask("task2")
	task3 := NewSuccessTask("task3")
	task4 := NewSuccessTask("task4")
	//task2.attributes[noop.AttrNameRetries] = "1"
	LinkDown(task1, task2)
	LinkDown(task1, task3)

	g := NewGroup()
	g.name = "group"

	g.AddNode(task1)
	g.AddNode(task2)
	g.AddNode(task3)
	g.AddNode(task4)
	g.Resolve()

	s := NewSession(g)
	g.AssignSession(s)
	r.RegisterSession(s)

	StartSessionForSuccess(s, false, t)

	assert.Equal(t, anemos.Fail, s.graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Fail, s.graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Fail, s.graph.nodes["task3"].Status())
	assert.Equal(t, anemos.Fail, s.graph.Status())
}

//func TestStress(t *testing.T) {
//
//	r := router.NewInternalRouter()
//
//	//channel := make(chan *Task)
//	//end := make(chan bool)
//	//
//	//go PrintChannel(channel)
//
//	depth := 1000
//	hor := make([][]*TaskNode, depth*2)
//	for i := 0; i < depth; i++ {
//		hor[i] = make([]*TaskNode, i+1)
//		hor[depth*2-1-i] = make([]*TaskNode, i+1)
//	}
//
//	g := NewGroup()
//	g.Name = "group"
//
//	hor[0][0] = NewSuccessTask("(0,0)")
//	g.AddNode(hor[0][0])
//
//	for i := 1; i < depth; i++ {
//		for j := 0; j <= i; j++ {
//			name := fmt.Sprintf("(%d,%d)", i, j)
//			//fmt.Println(name)
//
//			node := NewSuccessTask(name)
//			if j < i {
//				LinkDown(hor[i-1][j], node)
//			}
//			if j > 0 {
//				LinkDown(hor[i-1][j-1], node)
//			}
//			hor[i][j] = node
//			g.AddNode(node)
//
//		}
//	}
//
//	g.Resolve()
//	r.RegisterGroup(g)
//	StartSessionForSuccess(g, t)
//}
