package cfg

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zzell/transfer/db"
)

type (
	// Config configuration model
	Config struct {
		ListenPort        int            `json:"listen_port"`
		CommissionPercent float64        `json:"commission_percent"` // system commission
		Database          db.PostgresDSN `json:"db"`
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
