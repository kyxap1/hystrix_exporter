package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Cluster is a turbine cluster or a hystrix endpoint
type Cluster struct {
	Name string `yaml:"name,omitempty"`
	URL  string `yaml:"url,omitempty"`
}

// Config holds all clusters
type Config struct {
	Clusters []Cluster `yaml:"clusters,omitempty"`
}

// Parse a file into a config instance
func Parse(f string) (config Config, err error) {
	bts, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bts, &config)
	return
}
