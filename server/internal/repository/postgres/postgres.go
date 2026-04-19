package postgres

import (
	"fmt"
	"log"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func DBConnection() (*sqlx.DB, error) {
	cfg := config.Load()

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_NAME, cfg.DB_PASSWORD)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println(err)
		return nil, err
	}
	return db, err
}
