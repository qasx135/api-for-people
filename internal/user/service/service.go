package service

import (
	"api-for-people/internal/user/model"
	"context"
)

type RepositoryPostgres interface {
	Create(ctx context.Context, person *model.Person) error
	Get(ctx context.Context, id int) (model.Person, error)
	GetAll(ctx context.Context, params model.UserQueryParams) ([]model.Person, error)
	Update(ctx context.Context, person model.Person, id int) error
	Delete(ctx context.Context, id int) error
}

type Service struct {
	repo RepositoryPostgres
}

func NewService(repo RepositoryPostgres) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, person *model.Person) error {
	return s.repo.Create(ctx, person)
}
func (s *Service) Get(ctx context.Context, id int) (model.Person, error) {
	return s.repo.Get(ctx, id)
}
func (s *Service) GetAll(ctx context.Context, params model.UserQueryParams) ([]model.Person, error) {
	return s.repo.GetAll(ctx, params)
}
func (s *Service) Update(ctx context.Context, person model.Person, id int) error {
	return s.repo.Update(ctx, person, id)
}
func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
