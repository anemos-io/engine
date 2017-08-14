package clock

import (
	"testing"
	"time"
	//"fmt"
	//"sync"
	"github.com/stretchr/testify/assert"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
)

func Test_Clock(t *testing.T) {

	channel := make(chan *api.Event)
	clock := ClockObserver{
		channel:channel,
	}

	clock.Add("", time.Now().Add(time.Duration(3 * time.Second)))
	clock.Start()


	event := <- channel

	assert.Equal(t, "anemos/event:clock:0000:trigger", event.Uri)


	clock.Stop()
}
