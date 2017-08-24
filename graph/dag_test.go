package graph

import (
	"testing"
	//"fmt"
	//"sync"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func Equavalent_Dags(name string, t *testing.T) {
	group := ParseDagFile(fmt.Sprintf("%s.yaml", name))

	assert.Equal(t, 1, len(group.begin.downstream))
	assert.Equal(t, 2, len(group.end.upstream))

	taskStart := group.begin.downstream[fmt.Sprintf("$%s+begin>task-start", name)]
	assert.Equal(t, "task-start", taskStart.Name())
	assert.Equal(t, "anemos", taskStart.Provider())
	assert.Equal(t, "noop", taskStart.Operation())
	assert.Contains(t, taskStart.Downstream(), "task-start>task-left")
	taskLeft := taskStart.Downstream()["task-start>task-left"]
	assert.Equal(t, "task-left", taskLeft.Name())
	assert.Equal(t, "anemos", taskLeft.Provider())
	assert.Equal(t, "noop", taskLeft.Operation())
	assert.Contains(t, taskStart.Downstream(), "task-start>task-right")
	taskRight := taskStart.Downstream()["task-start>task-right"]
	assert.Equal(t, "task-right", taskRight.Name())
	assert.Equal(t, "anemos", taskRight.Provider())
	assert.Equal(t, "noop", taskRight.Operation())
}

func TestParseDag_ThreeTaskSimple(t *testing.T) {
	Equavalent_Dags("three-task-simple", t)
}

func TestParseDag_ThreeTaskFlat(t *testing.T) {
	Equavalent_Dags("three-task-flat", t)
}

func TestParseDag_ThreeTaskBinding(t *testing.T) {
	Equavalent_Dags("three-task-binding", t)
}
