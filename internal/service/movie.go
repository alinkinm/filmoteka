package service

import (
	"context"
	"filmoteka/internal/core"
	"fmt"
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
	SearchMovie(ctx context.Context, search string) (map[*core.Movie][]string, error)
}

type MovieService struct {
	movieRepository MovieRepository
}

func NewSegmentService(movieRepository MovieRepository, localRepository MovieRepository, minioRepository MovieRepository) *MovieService {
	return &MovieService{movieRepository: movieRepository}
}

func (service *MovieService) CreateMovie(ctx context.Context, movie *core.Movie, actors []int) error {
	return service.movieRepository.CreateMovie(ctx, movie, actors)
}

func (service *MovieService) DeleteMovie() error {
	return service.movieRepository.DeleteMovie()
}

func (service *MovieService) UpdateMovie(ctx context.Context, id int, columnName string, newValue interface{}) error {
	return service.movieRepository.UpdateMovie(ctx, id, columnName, newValue)
}

func (service *MovieService) AddActors(ctx context.Context, id int, actors []string) error {
	return service.movieRepository.AddActors(ctx, id, actors)
}

func (service *MovieService) DeleteActors(ctx context.Context, id int, actors []string) error {
	return service.movieRepository.DeleteActors(ctx, id, actors)
}

func (service *MovieService) GetAll(ctx context.Context, id int, actors []string, sorting string) (map[*core.Movie][]string, error) {
	switch sorting {
	case "rating":
		return service.movieRepository.GetAllMoviesByRating(ctx)
	case "title":
		return service.movieRepository.GetAllMoviesByTitle(ctx)
	case "release":
		return service.movieRepository.GetAllMoviesByReleaseDate(ctx)
	}
	return nil, fmt.Errorf("Internal server error")
}

func (service *MovieService) SearchMovie(ctx context.Context, search string) (map[*core.Movie][]string, error) {
	return service.movieRepository.SearchMovie(ctx, search)
}
