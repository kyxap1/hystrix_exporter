package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Cluster struct {
	Name string `yaml:"name,omitempty"`
	URL  string `yaml:"url,omitempty"`
}

type Config struct {
	Clusters []Cluster `yaml:"clusters,omitempty"`
}

func Parse(f string) (config Config, err error) {
	bts, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bts, &config)
	return
}
