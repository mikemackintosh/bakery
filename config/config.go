package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var Registry *Configuration

type Configuration struct {
	TempDir string `json:"tmp_dir" yaml:"tmp_dir"`
}

// NewFromFile loads the config at the filepath of f
func NewFromFile(f string) error {
	var ok bool
	var c *Configuration

	configFile, err := ioutil.ReadFile(f)
	if err != nil {
		return fmt.Errorf("Error while loading config from file %s", err)
	}

	// This should only really be used locally now that the configs are made and stored in datastore
	c, ok = isYaml(configFile)
	if !ok {
		c, ok = isJSON(configFile)
		if !ok {
			return fmt.Errorf("Configuration file is not in valid JSON or YAML format")
		}
	}
	Registry = c
	return nil
}

// isJSON tests if the file is JSON and returns decoded config
func isJSON(file []byte) (*Configuration, bool) {
	test := &Configuration{}
	if err := json.Unmarshal(file, &test); err != nil {
		return test, false
	}
	return test, true
}

// isYaml tests if the file is YAML and returns decoded config
func isYaml(file []byte) (*Configuration, bool) {
	test := &Configuration{}
	if err := yaml.Unmarshal(file, &test); err != nil {
		return test, false
	}
	return test, true
}
