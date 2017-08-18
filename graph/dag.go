package graph

import (
	//"fmt"
	//"github.com/ghodss/yaml"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	"log"
	"fmt"
	"github.com/anemos-io/engine"
)

type TaskConfig struct {
	Name       string
	TaskRef    string
	Downstream []TaskConfig
}

type BindingConfig struct {
}

type DagMetaData struct {
	Name string
}

type DagConfig struct {
	Kind     string
	Version  string
	MetaData DagMetaData
	Tasks    []TaskConfig
}

type builder struct {
	nodes    map[string]anemos.Node
	bindings map[string]binding
}
type binding struct {
	name string
	down string
	up   string
}

func parseTask(builder *builder, config map[interface{}]interface{}, parent anemos.Node) (anemos.Node) {
	name := config["name"].(string)

	task := NewTaskNode()
	builder.nodes[name] = task
	task.Name = name

	download, found := config["downstream"]
	if found {
		for _, v := range download.([]interface{}) {
			t := v.(map[interface{}]interface{})
			taskConfig, found := t["task"].(map[interface{}]interface{})
			if found {
				dt := parseTask(builder, taskConfig, task)
				//fmt.Println(task)
				bn := fmt.Sprintf("%s>%s", task.ShortName(), dt.ShortName())
				builder.bindings[bn] = binding{
					name: bn,
					up:   task.ShortName(),
					down: dt.ShortName(),
				}
			} else {
				taskRef, found := t["taskRef"].(string)
				if found {
					bn := fmt.Sprintf("%s>%s", task.ShortName(), taskRef)
					builder.bindings[bn] = binding{
						name: bn,
						up:   task.ShortName(),
						down: taskRef,
					}
				}
			}
		}
	}
	return task

}

func ParseDag(filename string) *Group {
	data, err := ioutil.ReadFile(filename)
	dagConfig := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &dagConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	builder := &builder{
		nodes:    make(map[string]anemos.Node),
		bindings: make(map[string]binding),
	}

	metaData := dagConfig["metaData"].(map[interface{}]interface{})

	//
	group := NewGroup()
	group.Name = metaData["name"].(string)

	tasksConfig := dagConfig["tasks"].([]interface{})
	fmt.Println(tasksConfig)
	for _, v := range tasksConfig {
		parseTask(builder, v.(map[interface{}]interface{}), nil)
	}

	c, found := dagConfig["bindings"]
	if found {
		bindingsConfig := c.([]interface{})
		fmt.Println(bindingsConfig)
		for _, v := range bindingsConfig {
			bindingConfig := v.(map[interface{}]interface{})
			n1 := bindingConfig["taskRef"].(string)
			d, found := bindingConfig["downstream"]
			if found {
				downstream := d.([]interface{})
				for _, d := range downstream {
					n2 := d.(map[interface{}]interface{})["taskRef"].(string)
					bn := fmt.Sprintf("%s>%s", n1, n2)
					builder.bindings[bn] = binding{
						name: bn,
						up:   n1,
						down: n2,
					}
				}
			}
			u, found := bindingConfig["upstream"]
			if found {
				upstream := u.([]interface{})
				for _, u := range upstream {
					n2 := u.(map[interface{}]interface{})["taskRef"].(string)
					bn := fmt.Sprintf("%s>%s", n2, n1)
					builder.bindings[bn] = binding{
						name: bn,
						up:   n2,
						down: n1,
					}
				}
			}

		}
	}


	for _, v := range builder.nodes {
		group.AddNode(v)
	}
	for _, v := range builder.bindings {
		LinkDownNamed(builder.nodes[v.up], builder.nodes[v.down], v.name)
	}
	group.Resolve()

	return group
}
