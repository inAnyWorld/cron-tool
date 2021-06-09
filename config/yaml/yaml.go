package yaml

import (
	"cacheflusher/config"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

var Yaml config.YamlConfig

func Loading() {
	yamlFile, yamlErr := ioutil.ReadFile("config/yaml/config.yaml")
	if yamlErr != nil {
		fmt.Printf("yaml 读取失败, Err:[%v]", yamlErr)
	}
	unmarshalErr := yaml.Unmarshal(yamlFile, &Yaml)
	if unmarshalErr != nil {
		fmt.Printf("yaml Unmarshal Err: [%v]", unmarshalErr)
	}
}
