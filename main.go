package main

import (
	"log"

	"github.com/zzell/transfer/cfg"
	"github.com/zzell/transfer/db"
	"github.com/zzell/transfer/db/repo"
)

const configFile = "config.json"

func main() {

	conf, err := cfg.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.NewDriver(&conf.Database)
	if err != nil {
		log.Fatal(err)
	}

	r := repo.NewWalletsRepo(conn)
	r.Transfer(1, 2, 500)
}
