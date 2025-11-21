package main

import (
	"log"

	"github.com/NH-Homelab/portfolio-backend/internal/config"
	pgdb "github.com/NH-Homelab/portfolio-backend/internal/pg_db"
)

func main() {
	backend_config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	dbconf := pgdb.Pg_Config{
		Host:     backend_config.Db_host,
		Port:     backend_config.Db_port,
		User:     backend_config.Db_user,
		Password: backend_config.Db_password,
		Db_name:  backend_config.Db_name,
	}

	pgdb, err := pgdb.NewPostgresDB(dbconf)
	if err != nil {
		log.Fatalf("Failed initial database setup: %v", err)
	}
	defer pgdb.Close()
}
