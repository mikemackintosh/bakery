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

func init() {
	Registry = &Configuration{}
}

// NewFromFile loads the config at the filepath of f
func NewFromFile(f string) error {

	configFile, err := ioutil.ReadFile(f)
	if err != nil {
		return fmt.Errorf("Error while loading config from file %s", err)
	}

	c, err := ParseConfig(configFile)
	if err != nil {
		return err
	}

	Registry = c

	return nil
}

// ParseConfig parses a configuration file
func ParseConfig(config []byte) (*Configuration, error) {
	var c *Configuration
	var ok bool
	// This should only really be used locally now that the configs are made and stored in datastore
	c, ok = isYaml(config)
	if !ok {
		c, ok = isJSON(config)
		if !ok {
			return c, fmt.Errorf("Configuration file is not in valid JSON or YAML format")
		}
	}
	return c, nil
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
