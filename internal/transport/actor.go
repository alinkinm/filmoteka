package transport

import (
	"context"
	"filmoteka/internal/core"
	"net/http"
)

type ActorService interface {
	CreateActor(ctx context.Context, actor *core.Actor) error
	UpdateActor(ctx context.Context, id int, columnName string, newValue interface{}) error
	DeleteActor(ctx context.Context, id int) error
	GetAllActors(ctx context.Context) ([]*core.Actor, error)
}

type ActorHandler struct {
	actorService ActorService
}

func NewActorHandler(service ActorService) *ActorHandler {
	return &ActorHandler{actorService: service}
}

func (handler *ActorHandler) CreateActor(w http.ResponseWriter, r *http.Request) {

}
