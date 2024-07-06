package postgres

import (
	"database/sql"
	"fmt"
	"github.com/RomanLevBy/BurgersAPI/internal/config"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(conf *config.Config) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			conf.Postgres.Host,
			conf.Postgres.User,
			conf.Postgres.DBName,
			conf.Postgres.Password,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	return &Storage{db: db}, nil
}
