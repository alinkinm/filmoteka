package service

import (
	"context"
	"filmoteka/internal/core"
)

type ActorRepository interface {
	CreateActor(ctx context.Context, actor *core.Actor) error
	UpdateActor(ctx context.Context, id int, columnName string, newValue interface{}) error
	DeleteActor(ctx context.Context, id int) error
	GetAllActors(ctx context.Context) ([]*core.Actor, error)
}

type ActorService struct {
	actorRepository ActorRepository
}

func NewActorService(actorRepository ActorRepository) *ActorService {
	return &ActorService{actorRepository: actorRepository}
}

func (service *ActorService) CreateActor(ctx context.Context, actor *core.Actor) error {
	return service.actorRepository.CreateActor(ctx, actor)
}

func (service *ActorService) UpdateActor(ctx context.Context, id int, columnName string, newValue interface{}) error {
	return service.actorRepository.UpdateActor(ctx, id, columnName, newValue)
}

func (service *ActorService) DeleteActor(ctx context.Context, id int) error {
	return service.actorRepository.DeleteActor(ctx, id)
}

func (service *ActorService) GetAllActors(ctx context.Context) ([]*core.Actor, error) {
	return service.actorRepository.GetAllActors(ctx)
}
