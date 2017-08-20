package store

import (
	"fmt"
	"github.com/anemos-io/engine/graph"
	"io/ioutil"

	"log"
	"os"
)

func NewGraphStore() *GraphStore {
	gs := &GraphStore{
		graphs: make(map[string]*graph.Group),
	}
	gs.LoadAll()

	log.Printf("Loaded %d graphs", len(gs.graphs))
	return gs
}

type GraphStore struct {
	graphs map[string]*graph.Group
}

func (g *GraphStore) LoadAll() {

	home := os.Getenv("HOME")
	base := ".anemos/graphs"
	dirs, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s", home, base))
	for _, dir := range dirs {
		fn := fmt.Sprintf("%s/%s/%s", home, base, dir.Name())
		group := graph.ParseDagFile(fn)
		g.graphs[group.Name()] = group
	}

}
