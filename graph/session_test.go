package graph

import (
	"github.com/anemos-io/engine"
	"github.com/anemos-io/engine/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopySessionShouldNotAffectOriginal(t *testing.T) {

	r := router.NewRouter(nil)

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

	assert.Equal(t, anemos.Unknown, task1.Status())
	assert.Equal(t, anemos.Unknown, task2.Status())
	assert.Equal(t, anemos.Unknown, g.Status())

	assert.Equal(t, anemos.Success, s.Graph.nodes["task1"].Status())
	assert.Equal(t, anemos.Success, s.Graph.nodes["task2"].Status())
	assert.Equal(t, anemos.Success, s.Graph.Status())
}
