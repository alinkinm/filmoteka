package repository

import (
	"context"
	"filmoteka/internal/core"
	"fmt"

	log "github.com/sirupsen/logrus"

	"errors"

	"database/sql"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type ActorRepository struct {
	Db *sqlx.DB
}

const (
	CreateActor           = "INSERT INTO Actors(names, sex, bd) SELECT $1, $2, $3 returning id;"
	DeleteActor           = "DELETE FROM Actors where id = ($1) returning id;"
	DeleteActorFromMovies = "DELETE FROM ActorMovie where actor_id = $1;"
	UpdateActor           = "UPDATE Actors SET $1 = $2 where id = $3;"
	GetAllActors          = `SELECT a.id AS actor_id, a.names AS actor_name, a.sex AS actor_sex, a.names AS actor_name, a.bd AS actor_birthday, 
	array_agg(m.title) AS movies_participated FROM Actors a LEFT JOIN ActorMovie am ON a.id = am.actor_id LEFT JOIN Movies m ON am.movie_id = m.id GROUP BY a.id, a.names;`
)

func NewActorRepository(db *sqlx.DB) *ActorRepository {
	return &ActorRepository{Db: db}
}

func (repository *ActorRepository) CreateActor(ctx context.Context, actor *core.Actor) error {

	_, err := repository.Db.ExecContext(ctx, CreateActor, actor.Name, string(actor.Sex), actor.Bd)

	var e *pgconn.PgError
	if err != nil {
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			log.Info(err.Error())
			return core.NewErrActorAlreadyExists()
		}
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	return nil

}

func (repository *ActorRepository) UpdateActor(ctx context.Context, id int, columnName string, newValue interface{}) error {

	res, err := repository.Db.ExecContext(ctx, UpdateActor, columnName, newValue, id)

	if err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	if rows == 0 {
		log.Info("No rows affected")
		return core.NewErrActorDoesNotExist()
	}

	return nil
}

func (repository *ActorRepository) DeleteActor(ctx context.Context, id int) error {

	tx, err := repository.Db.Begin()
	if err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, DeleteActor, id)

	if err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	if rows == 0 {
		return core.NewErrActorDoesNotExist()
	}

	_, err = tx.ExecContext(ctx, DeleteActorFromMovies, id)

	if err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	if err = tx.Commit(); err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	return nil

}

func (repository *ActorRepository) GetAllActors(ctx context.Context) (map[*core.Actor][]string, error) {
	var actors map[*core.Actor][]string

	rows, err := repository.Db.QueryContext(ctx, GetAllActors)
	if err == sql.ErrNoRows {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}
	if err != nil {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}

	for rows.Next() {
		actor := &core.Actor{}
		actors_movies := []string{}
		err = rows.Scan(&actor.Id, &actor.Name, &actor.Sex, &actor.Bd, actors_movies)

		if err != nil {
			log.Info(err.Error())
			return nil, fmt.Errorf("Internal server error")
		}
		actors[actor] = actors_movies
	}

	return actors, nil
}
