package sourceapp

import (
	"context"
	"fmt"
	"sync"
	"time"

	domain "github.com/acme/streamforge-cdc/internal/source/domain"
)

type Repository interface {
	Save(context.Context, domain.Source) error
	Get(context.Context, string) (domain.Source, error)
	List(context.Context) ([]domain.Source, error)
}
type Validator interface {
	Validate(context.Context, domain.Source) error
}
type Service struct {
	repo      Repository
	validator Validator
	mu        sync.Mutex
}

func New(r Repository, v Validator) *Service { return &Service{repo: r, validator: v} }
func (s *Service) Create(ctx context.Context, src domain.Source) (domain.Source, error) {
	if src.ID == "" {
		return src, fmt.Errorf("source id is required")
	}
	if src.Kind != domain.PostgreSQL && src.Kind != domain.MySQL {
		return src, fmt.Errorf("unsupported source kind %q", src.Kind)
	}
	now := time.Now().UTC()
	src.Status = domain.Draft
	src.CreatedAt = now
	src.UpdatedAt = now
	if err := s.repo.Save(ctx, src); err != nil {
		return src, fmt.Errorf("create source: %w", err)
	}
	return src, nil
}
func (s *Service) Activate(ctx context.Context, id string) error {
	src, err := s.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("activate: %w", err)
	}
	src.Status = domain.Validating
	src.UpdatedAt = time.Now().UTC()
	_ = s.repo.Save(ctx, src)
	if s.validator != nil {
		if err = s.validator.Validate(ctx, src); err != nil {
			src.Status = domain.Failed
			_ = s.repo.Save(ctx, src)
			return fmt.Errorf("validate source: %w", err)
		}
	}
	src.Status = domain.Active
	src.UpdatedAt = time.Now().UTC()
	if err = s.repo.Save(ctx, src); err != nil {
		return fmt.Errorf("activate save: %w", err)
	}
	return nil
}
func (s *Service) SetStatus(ctx context.Context, id string, status domain.Status) error {
	src, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	src.Status = status
	src.UpdatedAt = time.Now().UTC()
	if err = s.repo.Save(ctx, src); err != nil {
		return fmt.Errorf("status: %w", err)
	}
	return nil
}
func (s *Service) List(ctx context.Context) ([]domain.Source, error) { return s.repo.List(ctx) }
