package infrastructure

import (
	"context"
	"fmt"
	"log"

	"filmoteka/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

func SetUpPostgresDatabase(ctx context.Context, config *config.DBConfig) (*sqlx.DB, error) {

	db, err := sqlx.ConnectContext(ctx, "pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		config.Host, config.Port, config.Name, config.Password, config.Db))

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return db, nil

}
