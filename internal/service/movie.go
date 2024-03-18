package service

import (
	"context"
	"filmoteka/internal/core"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *core.Movie, actors []int) error
	DeleteMovie() error
	UpdateMovie(ctx context.Context, id int, columnName string, newValue interface{}) error
	AddActors(ctx context.Context, id int, actors []string) error
	DeleteActors(ctx context.Context, id int, actors []string) error
	GetAllMoviesByRating(ctx context.Context) (map[*core.Movie][]string, error)
	GetAllMoviesByTitle(ctx context.Context) (map[*core.Movie][]string, error)
	GetAllMoviesByReleaseDate(ctx context.Context) (map[*core.Movie][]string, error)
	SearchMovieByTitle(ctx context.Context) (map[*core.Movie][]string, error)
	SearchMovieByActor(ctx context.Context) (map[*core.Movie][]string, error)
}

type MovieService struct {
	movieRepository MovieRepository
}

func NewSegmentService(movieRepository MovieRepository, localRepository MovieRepository, minioRepository MovieRepository) *MovieService {
	return &MovieService{movieRepository: movieRepository}
}

func 
