package graph

import (
	"fmt"
	"github.com/anemos-io/engine"
	api "github.com/anemos-io/engine/grpc/anemos/v1alpha1"
	"math/rand"
)

type Session struct {
	router anemos.Router
	graph  *Group

	id int

	Instances map[anemos.Node][]*api.TaskInstance
}

func (session *Session) SetRouter(router anemos.Router) {
	session.router = router
}

func (session *Session) Router() anemos.Router {
	return session.router
}

func (session *Session) NewTaskInstance(n anemos.Node) *api.TaskInstance {
	def := api.TaskInstance{
		Provider:   n.Provider(),
		Operation:  n.Operation(),
		Name:       n.Name(),
		Id:         fmt.Sprintf("%x/%d/%d", rand.Int63(), session.id, 1),
		Attributes: make(map[string]string, 0),
		Metadata:   make(map[string]string, 0),
	}
	for key, value := range n.Attributes() {
		def.Attributes[fmt.Sprintf("anemos/attribute:%s:%s:%s", n.Provider(), n.Operation(), key)] = value
	}
	instances, found := session.Instances[n]
	if !found {
		instances = make([]*api.TaskInstance, 1)
		session.Instances[n] = instances
	}
	instances = append(instances, &def)
	return &def
}

func NewSession(original *Group) *Session {

	session := &Session{
		Instances: make(map[anemos.Node][]*api.TaskInstance, 0),
	}

	return session
}
