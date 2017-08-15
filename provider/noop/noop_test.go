package noop

import (
	"testing"
	"github.com/stretchr/testify/assert"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"github.com/anemos-io/engine"
)

func execute(instance *api.TaskInstance) (*api.Event) {

	channel := make(chan *api.Event)

	executor := NoopExecutor{}
	observer := NoopObserver{
		EventChannel: channel,
	}

	executor.CoupleObserver(&observer)
	executor.Execute(instance)

	return <-channel
}

func newInstance() (*api.TaskInstance) {
	return &api.TaskInstance{
		Name:       "test",
		Id:         "0042",
		Attributes: make(map[string]string),
		Metadata:   make(map[string]string),
	}

}

func TestNoop_ExplicitZeroRetry(t *testing.T) {

	instance := newInstance()

	instance.Metadata[anemos.MetaTaskRetry] = "0"
	instance.Attributes[AttrSuccessDuration] = "1ms"
	instance.Attributes[AttrRetries] = "0"
	event := execute(instance)

	uri, _ := anemos.ParseUri(event.Uri)
	assert.Equal(t, "success", uri.Status)
}

func TestNoop_AllDefault(t *testing.T) {

	instance := newInstance()
	event := execute(instance)

	uri, _ := anemos.ParseUri(event.Uri)
	assert.Equal(t, "success", uri.Status)
}

func TestNoop_FailOnRetries1Retry0(t *testing.T) {

	instance := newInstance()

	instance.Metadata[anemos.MetaTaskRetry] = "0"
	instance.Attributes[AttrRetries] = "1"
	event := execute(instance)

	uri, _ := anemos.ParseUri(event.Uri)
	assert.Equal(t, "fail", uri.Status)
}

func TestNoop_SucceedOnRetries1Retry1(t *testing.T) {

	instance := newInstance()

	instance.Metadata[anemos.MetaTaskRetry] = "1"
	instance.Attributes[AttrRetries] = "1"
	event := execute(instance)

	uri, _ := anemos.ParseUri(event.Uri)
	assert.Equal(t, "success", uri.Status)
}
