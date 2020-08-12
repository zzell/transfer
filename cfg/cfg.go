package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type (
	// Config configuration model
	Config struct {
		Database DB `json:"db"`
	}

	// DB configuration model
	DB struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Server   string `json:"server"`
		Port     string `json:"port"`
		Database string `json:"database"`
	}
)

// New reads config from file
func New(file string) (c Config, err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	return c, json.Unmarshal(b, &c)
}
