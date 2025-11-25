package main

import (
	"log"
	"net/http"

	"github.com/NH-Homelab/portfolio-backend/internal/config"
	pgdb "github.com/NH-Homelab/portfolio-backend/internal/pg_db"
	portfoliodao "github.com/NH-Homelab/portfolio-backend/internal/portfolio_dao"
	publichandler "github.com/NH-Homelab/portfolio-backend/internal/public_handler"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s%s from %s\n", r.Method, r.Host, r.URL, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

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

	dao := portfoliodao.NewPortfolioDao(pgdb)
	ph := publichandler.NewPublicHandler(dao)
	mux := http.NewServeMux()

	ph.RegisterHandlers(mux)

	log.Printf("Starting HTTP server on :8080...")
	err = http.ListenAndServe(":8080", logRequest(setContentType(mux)))
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
