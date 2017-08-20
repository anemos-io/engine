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
		Graphs: make(map[string]*graph.Group),
	}
	gs.LoadAll()

	log.Printf("Loaded %d Graphs", len(gs.Graphs))
	return gs
}

type GraphStore struct {
	Graphs map[string]*graph.Group
}

func (g *GraphStore) LoadAll() {

	home := os.Getenv("HOME")
	base := ".anemos/Graphs"
	dirs, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s", home, base))
	for _, dir := range dirs {
		fn := fmt.Sprintf("%s/%s/%s", home, base, dir.Name())
		group := graph.ParseDagFile(fn)
		g.Graphs[group.Name()] = group
	}

}
