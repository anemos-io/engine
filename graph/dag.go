package graph

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type TaskConfig struct {
	Name string
	TaskRef string
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

func ParseDag(filename string) *DagConfig {
	dat, err := ioutil.ReadFile(filename)
	var dagConfig DagConfig
	err = yaml.Unmarshal(dat, &dagConfig)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	return &dagConfig

	//fmt.Println(dagConfig.Tasks[0].Name)
	//fmt.Println(dagConfig.Tasks[0].Body) //["reference"])
	//
	//for k, v := range dagConfig.Tasks[0].Body {
	//	fmt.Printf("key[%s] value[%s]\n", k, v)
	//}

}
