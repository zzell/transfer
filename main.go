package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zzell/transfer/cfg"
	"github.com/zzell/transfer/db"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/web"
)

const configFile = "config.json"

func main() {
	config, err := cfg.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.NewDriver(&config.Database)
	if err != nil {
		log.Fatal(err)
	}

	walletsRepo := repo.NewWalletsRepo(conn)
	router := web.NewRouter(walletsRepo, &config)
	fmt.Printf("listening on port :%d\n", config.ListenPort)
	log.Print(http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), router))
}
