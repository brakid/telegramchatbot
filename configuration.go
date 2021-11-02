package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const CONFIG_PATH = "bot.cfg"

type Configuration struct {
	configurations map[string]string
}

func (c *Configuration) Set(key, value string) {
	(*c).configurations[key] = value
}

func (c *Configuration) Get(key string) (string, bool) {
	value, ok := (*c).configurations[key]
	return value, ok
}

func (c *Configuration) String() string {
	return fmt.Sprintf("Configuration: %s", c.configurations)
}

func LoadConfiguration() (*Configuration, error) {
	file, err := os.Open(CONFIG_PATH)
	if os.IsNotExist(err) {
		return &Configuration{configurations: make(map[string]string)}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configurations := make(map[string]string)

	err = json.NewDecoder(file).Decode(&configurations)
	if err != nil {
		return nil, err
	}

	return &Configuration{configurations: configurations}, err
}

func (c *Configuration) Save() error {
	marshalledConfiguration, err := json.MarshalIndent(c.configurations, "", "\t")
	if err != nil {
		return err
	}

	reader := bytes.NewReader(marshalledConfiguration)

	file, err := os.Create(CONFIG_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}
