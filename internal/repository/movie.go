package repository

import (
	"context"
	"errors"
	"filmoteka/internal/core"
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"strings"

	"database/sql"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type MovieRepository struct {
	Db *sqlx.DB
}

const (
	UniqueViolationErr  = "23505"
	ForeignKeyViolation = "23503"

	CreateMovie = "INSERT INTO Movies(title, descr, release, rating) SELECT $1, $2, $3, $4 returning id;"
	UpdateMovie = "UPDATE movies SET $1 = $2 where id = $3;"

	AddActorsToMovie      = "SELECT add_actors_to_movie($1, $2);"
	DeleteActorsFromMovie = "SELECT delete_actors_from_movie($1, $2);"
	DeleteMovie           = "DELETE FROM Actors where id = ($1) returning *;"

	SortMoviesByRating = ` SELECT id, title, descr, release, rating, array_agg(a.names) AS actors FROM Movies m LEFT JOIN ActorMovie am ON m.id = am.movie_id 
	LEFT JOIN Actors a ON am.actor_id = a.id GROUP BY m.id, m.title, m.descr, m.release, m.rating ORDER BY m.rating DESC;`

	SortMoviesByReleaseDate = ` SELECT id, title, descr, release, rating, array_agg(a.names) AS actors FROM Movies m LEFT JOIN ActorMovie am ON m.id = am.movie_id
	LEFT JOIN Actors a ON am.actor_id = a.id GROUP BY m.id, m.title, m.descr, m.release, m.rating ORDER BY m.release DESC;`

	SortMoviesByTitle = `SELECT id, title, descr, release, rating, array_agg(a.names) AS actors FROM Movies m LEFT JOIN ActorMovie am ON m.id = am.movie_id 
	LEFT JOIN Actors a ON am.actor_id = a.id GROUP BY m.id, m.title, m.descr, m.release, m.rating ORDER BY m.title;`

	SearchMovie = `SELECT m.id AS movie_id, m.title AS movie_title, m.descr, m.release, m.rating, array_agg(a.names) AS actors FROM Movies m
	LEFT JOIN ActorMovie am ON m.id = am.movie_id LEFT JOIN Actors a ON am.actor_id = a.id WHERE m.title ILIKE '%' || $1 || '%' OR a.names ILIKE '%' || $1 || '%'
	GROUP BY m.id, m.title, m.descr, m.release, m.rating;`
)

func NewMovieRepository(db *sqlx.DB) *MovieRepository {
	return &MovieRepository{Db: db}
}

func (repository *MovieRepository) CreateMovie(ctx context.Context, movie *core.Movie) error {

	tx, err := repository.Db.Begin()
	if err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	defer tx.Rollback()

	newid := -1
	err = tx.QueryRowContext(ctx, CreateMovie, movie.Title, movie.Descr, movie.Release, movie.Rating).Scan(&newid)
	movie.Id = newid

	var e *pgconn.PgError
	if err != nil {
		log.Info(err.Error())
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return core.NewErrMovieAlreadyExists()
		}
		return fmt.Errorf("Internal server error")
	}

	_, err = tx.ExecContext(ctx, AddActorsToMovie, movie.Actors, movie.Id)

	if err != nil {
		log.Info(err.Error())
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			var nonExistentActorID int
			_, err := fmt.Sscanf(e.Detail, "Key (actor_id)=(%d)", &nonExistentActorID)
			newerr := core.NewErrActorDoesNotExist()
			newerr.Inf.Msg += fmt.Sprintf("Non-existent actor ID: %d\n", nonExistentActorID)
			log.Info(err.Error())
			return newerr
		}
		return fmt.Errorf("Internal server error")
	}

	if err = tx.Commit(); err != nil {
		log.Info(err.Error())
		return fmt.Errorf("Internal server error")
	}

	return nil

}

func (repository *MovieRepository) DeleteMovie() error {
	return nil
}

func (repository *MovieRepository) UpdateMovie(ctx context.Context, id int, columnName string, newValue interface{}) error {
	res, err := repository.Db.ExecContext(ctx, UpdateMovie, columnName, newValue, id)

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
		return core.NewErrMovieDoesNotExist()
	}

	return nil

}

func (repository *MovieRepository) AddActors(ctx context.Context, id int, actors []string) error {

	_, err := repository.Db.ExecContext(ctx, AddActorsToMovie, actors, id)

	if err != nil {
		var e *pgconn.PgError
		log.Info(err.Error())
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			detail := e.Detail
			var missingID int
			var missingType string
			if strings.Contains(detail, "actor_id") {
				missingType = "actor_id"
			} else if strings.Contains(detail, "movie_id") {
				missingType = "movie_id"
			}
			_, err := fmt.Sscanf(detail, "Key (%s)=(%d)", &missingType, &missingID)
			log.Info(err.Error())
			return fmt.Errorf(fmt.Sprintf("Missing ID: %d, Type: %s\n", missingID, missingType))
		}
		return fmt.Errorf("Internal server error")
	}

	return nil
}

func (repository *MovieRepository) DeleteActors(ctx context.Context, id int, actors []string) error {
	_, err := repository.Db.ExecContext(ctx, DeleteActorsFromMovie, actors, id)

	if err != nil {
		var e *pgconn.PgError
		log.Info(err.Error())
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			detail := e.Detail
			var missingID int
			var missingType string
			if strings.Contains(detail, "actor_id") {
				missingType = "actor_id"
			} else if strings.Contains(detail, "movie_id") {
				missingType = "movie_id"
			}
			_, err := fmt.Sscanf(detail, "Key (%s)=(%d)", &missingType, &missingID)
			log.Info(err.Error())
			return fmt.Errorf(fmt.Sprintf("Missing ID: %d, Type: %s\n", missingID, missingType))
		}
		return fmt.Errorf("Internal server error")
	}

	return nil
}

func (repository *MovieRepository) GetAllMoviesByRating(ctx context.Context) ([]*core.Movie, error) {

	var movies []*core.Movie

	rows, err := repository.Db.QueryContext(ctx, SortMoviesByRating)
	if err == sql.ErrNoRows {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}
	if err != nil {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}

	for rows.Next() {
		movie := &core.Movie{}
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Descr, &movie.Release, &movie.Rating, &movie.Actors)

		if err != nil {
			log.Info(err.Error())
			return nil, fmt.Errorf("Internal server error")
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (repository *MovieRepository) GetAllMoviesByTitle(ctx context.Context) ([]*core.Movie, error) {
	var movies []*core.Movie

	rows, err := repository.Db.QueryContext(ctx, SortMoviesByTitle)
	if err == sql.ErrNoRows {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}
	if err != nil {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}

	for rows.Next() {
		movie := &core.Movie{}
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Descr, &movie.Release, &movie.Rating, &movie.Actors)

		if err != nil {
			log.Info(err.Error())
			return nil, fmt.Errorf("Internal server error")
		}
		movies = append(movies, movie)
	}

	return movies, nil

}

func (repository *MovieRepository) GetAllMoviesByReleaseDate(ctx context.Context) ([]*core.Movie, error) {
	var movies []*core.Movie

	rows, err := repository.Db.QueryContext(ctx, SortMoviesByReleaseDate)
	if err == sql.ErrNoRows {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}
	if err != nil {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}

	for rows.Next() {
		movie := &core.Movie{}
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Descr, &movie.Release, &movie.Rating, &movie.Actors)

		if err != nil {
			log.Info(err.Error())
			return nil, fmt.Errorf("Internal server error")
		}
		movies = append(movies, movie)
	}

	return movies, nil

}

func (repository *MovieRepository) SearchMovie(ctx context.Context, search string) ([]*core.Movie, error) {
	var movies []*core.Movie

	rows, err := repository.Db.QueryContext(ctx, SearchMovie, search)
	if err == sql.ErrNoRows {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}
	if err != nil {
		log.Info(err.Error())
		return nil, fmt.Errorf("Internal server error")
	}

	for rows.Next() {
		movie := &core.Movie{}
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Descr, &movie.Release, &movie.Rating, &movie.Actors)

		if err != nil {
			log.Info(err.Error())
			return nil, fmt.Errorf("Internal server error")
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
