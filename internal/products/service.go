package products

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/luizflm/go-ecom-api/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, id int64) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductByID(ctx context.Context, id int64) (repo.Product, error) {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return repo.Product{}, ErrNotFound
		default:
			return repo.Product{}, err
		}
	}

	return product, nil
}
