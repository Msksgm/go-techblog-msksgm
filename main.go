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
		log.Fatalf("error is occuered because %v", err)
	}

	db, err := postgres.Open(cfg.dbURI)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	srv := server.NewServer(db)
	log.Fatal(srv.Run(cfg.port))
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
