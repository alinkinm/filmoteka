package transport

import (
	"context"
	"filmoteka/internal/core"
	"net/http"
	"time"
)

type MovieService interface {
	CreateMovie(ctx context.Context, movie *core.Movie, actors []int) error
	DeleteMovie() error
	UpdateMovie(ctx context.Context, id int, columnName string, newValue interface{}) error
	AddActors(ctx context.Context, id int, actors []string) error
	DeleteActors(ctx context.Context, id int, actors []string) error
	SearchMovie(ctx context.Context, search string) (map[*core.Movie][]string, error)
}

type MovieHandler struct {
	segmentService MovieService
}

func NewMovieHandler(service MovieService) *MovieHandler {
	return &MovieHandler{segmentService: service}
}

func (handler *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request, ctx context.Context) error {

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Millisecond)
	defer cancel()

	r = r.WithContext(ctx)

}
