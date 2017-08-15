package noop

import (
	"testing"
	//"time"
	//"fmt"
	//"sync"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"fmt"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"github.com/anemos-io/engine"
)

func Test_time(t *testing.T) {

	//val, err := time.Parse(time.RFC3339, "2017-11-28T15:15:30Z")
	//if err != nil { // Always check errors even if they should not happen.
	//	panic(err)
	//}
	//fmt.Println(val)

	channel := make(chan *api.Event)

	executor := NoopExecutor{}
	observer := NoopObserver{
		EventChannel: channel,
	}

	executor.CoupleObserver(&observer)

	instance := &api.TaskInstance{
		Name:       "test",
		Id:         "0042",
		Attributes: make(map[string]string),
		Metadata:   make(map[string]string),
	}

	instance.Metadata[anemos.MetaRetry] = "0"
	instance.Attributes[AttrSuccessDuration] = "1ms"
	instance.Attributes[AttrRetries] = "0"
	executor.Execute(instance)

	event := <-channel

	assert.Equal(t, "anemos/event:anemos:noop:test:0042:success", event.Uri)
	fmt.Println(event.Metadata)

	//dagConfig := ParseDag("three-task-simple.yaml")
	//
	//assert.Equal(t, "anemos/dag", dagConfig.Kind)
	//assert.Equal(t, "v1", dagConfig.Version)
	//assert.Equal(t, "three-node-simple", dagConfig.MetaData.Name)
	//assert.Equal(t, "task-start", dagConfig.Tasks[0].Name)
	//assert.Equal(t, "task-top", dagConfig.Tasks[0].Downstream[0].Name)
	//assert.Equal(t, "task-bottom", dagConfig.Tasks[0].Downstream[1].Name)
}
