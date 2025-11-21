package pgdb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Conn *sql.DB
}

type Pg_Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Db_name  string
}

// Instantiates a PostgresDB connection type
func NewPostgresDB(config Pg_Config) (*PostgresDB, error) {
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Db_name)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{Conn: db}, nil
}

func (pg *PostgresDB) Close() error {
	return pg.Conn.Close()
}

func (pg *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return pg.Conn.Exec(query, args...)
}

func (pg *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return pg.Conn.Query(query, args...)
}
