package graph

import (
	"testing"
	//"fmt"
	//"sync"
	"github.com/stretchr/testify/assert"
)

func TestParseDag_ThreeTaskSimple(t *testing.T) {
	dagConfig := ParseDag("three-task-simple.yaml")

	assert.Equal(t, "anemos/dag", dagConfig.Kind)
	assert.Equal(t, "v1", dagConfig.Version)
	assert.Equal(t, "three-node-simple", dagConfig.MetaData.Name)
	assert.Equal(t, "task-begin", dagConfig.Tasks[0].Name)
	assert.Equal(t, "task-top", dagConfig.Tasks[0].Downstream[0].Name)
	assert.Equal(t, "task-bottom", dagConfig.Tasks[0].Downstream[1].Name)
}

