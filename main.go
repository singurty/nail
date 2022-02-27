package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/joho/godotenv"

	"github.com/singurty/nail/db"
	"github.com/singurty/nail/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer db.DBpool.Close()

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
