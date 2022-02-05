package main

import (
	"fmt"
	"log"
	"os"

	"github.com/msksgm/go-techblog-msksgm/postgres"
	"github.com/msksgm/go-techblog-msksgm/server"
)

type config struct {
	port  string
	dbURI string
}

func main() {
	cfg, err := envConfig()
	if err != nil {
		log.Fatalln("err:", err)
	}

	db, err := postgres.Open(cfg.dbURI)
	if err != nil {
		log.Fatalln("err:", err)
	}

	srv := server.NewServer(db)
	if err := srv.Run(cfg.port); err != nil {
		log.Fatalln("err:", err)
	}
}

func envConfig() (config, error) {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return config{}, fmt.Errorf("PORT is not provided")
	}

	dbURI, ok := os.LookupEnv("POSTGRESQL_URL")
	if !ok {
		return config{}, fmt.Errorf("POSTGRESQL_URL is not provided")
	}

	return config{port: port, dbURI: dbURI}, nil
}
