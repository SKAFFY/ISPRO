package service

import (
	"context"
	"errors"

	"github.com/mkheyfets/ispro-app/models"
	"github.com/mkheyfets/ispro-app/pkg/repository"
)

var (
	ErrNotFound = errors.New("not found")
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListEntries(ctx context.Context) ([]*models.Entry, error) {
	return s.repo.ListEntries(ctx)
}

func (s *Service) GetEntry(ctx context.Context, id int64) (*models.Entry, error) {
	entry, err := s.repo.GetEntry(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return entry, nil
}

func (s *Service) CreateEntry(ctx context.Context, title, content string) (*models.Entry, error) {
	return s.repo.CreateEntry(ctx, title, content)
}

func (s *Service) UpdateEntry(ctx context.Context, id int64, title, content string) (*models.Entry, error) {
	return s.repo.UpdateEntry(ctx, id, title, content)
}

func (s *Service) DeleteEntry(ctx context.Context, id int64) error {
	return s.repo.DeleteEntry(ctx, id)
}

func (s *Service) ListLinks(ctx context.Context) ([]*models.Link, error) {
	return s.repo.ListLinks(ctx)
}

func (s *Service) GetLink(ctx context.Context, id int64) (*models.Link, error) {
	link, err := s.repo.GetLink(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return link, nil
}

func (s *Service) CreateLink(ctx context.Context, sourceID, targetID int64) (*models.Link, error) {
	return s.repo.CreateLink(ctx, sourceID, targetID)
}

func (s *Service) UpdateLink(ctx context.Context, id int64, sourceID, targetID int64) (*models.Link, error) {
	return s.repo.UpdateLink(ctx, id, sourceID, targetID)
}

func (s *Service) DeleteLink(ctx context.Context, id int64) error {
	return s.repo.DeleteLink(ctx, id)
}