package main

import (
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
	cfg := envConfig()

	db, err := postgres.Open(cfg.dbURI)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	srv := server.NewServer(db)
	log.Fatal(srv.Run(cfg.port))
}

func envConfig() config {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		panic("PORT not provided")
	}

	dbURI, ok := os.LookupEnv("POSTGRESQL_URL")
	if !ok {
		panic("POSTGRESQL_URL not provided")
	}

	return config{port: port, dbURI: dbURI}
}
