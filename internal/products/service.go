package products

import (
	"context"

	repo "github.com/ybuilds/ecomm-api/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	ListProductById(ctx context.Context, id int64) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (svc *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return svc.repo.ListProducts(ctx)
}

func (svc *svc) ListProductById(ctx context.Context, id int64) (repo.Product, error) {
	return svc.repo.FindProductByID(ctx, id)
}
