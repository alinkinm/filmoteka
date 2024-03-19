package transport

import (
	"context"
	"filmoteka/internal/core"
	"net/http"

	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type MovieService interface {
	CreateMovie(ctx context.Context, movie *core.Movie) error
	DeleteMovie() error
	UpdateMovie(ctx context.Context, id int, columnName string, newValue interface{}) error
	AddActors(ctx context.Context, id int, actors []string) error
	DeleteActors(ctx context.Context, id int, actors []string) error
	SearchMovie(ctx context.Context, search string) ([]*core.Movie, error)
}

type MovieHandler struct {
	movieService MovieService
}

func NewMovieHandler(service MovieService) *MovieHandler {
	return &MovieHandler{movieService: service}
}

type KeyMovie struct{}

func (handler *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {

	movie := &core.Movie{}
	log.Info("Parsing movie from request")

	d := json.NewDecoder(r.Body)
	err := d.Decode(&movie)

	if err != nil {
		log.Info("Could not parse movie from request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	log.Info("Parsed movie - ", movie)

	err = handler.movieService.CreateMovie(r.Context(), movie)
	if err != nil {
		log.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}
