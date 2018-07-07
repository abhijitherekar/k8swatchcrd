package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var ConfigFileName = ".k8swatchcrd.yaml"

// Resource contains resource configuration
type Resource struct {
	Deployment            bool `json:"deployment"`
	ReplicationController bool `json:"rc"`
	ReplicaSet            bool `json:"rs"`
	DaemonSet             bool `json:"ds"`
	Services              bool `json:"svc"`
	Pod                   bool `yaml:"pod"`
	Job                   bool `json:"job"`
	PersistentVolume      bool `json:"pv"`
}

// Config struct contains k8swatchcrd's configuration
type Config struct {
	//Reason   []string `json:"reason"`
	Resource Resource `json:"resource"`
}

// Load loads configuration from config file
func (c *Config) Load() error {

	file, err := os.Open(os.Getenv("HOME") + "/" + ConfigFileName)
	if err != nil {
		fmt.Println("\n could not open the file")
		return err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("\n error while reading")
		return err
	}
	if len(b) != 0 {
		return yaml.Unmarshal(b, c)
	}

	return nil
}

// New creates new config object
func New() (*Config, error) {
	c := &Config{}
	if err := c.Load(); err != nil {
		return c, err
	}

	return c, nil
}
