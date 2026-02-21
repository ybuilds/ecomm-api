package products

import "context"

type Service interface {
	ListProducts(ctx context.Context) error
}

type svc struct {
	//repository dependency
}

func NewService() Service {
	return &svc{}
}

func (svc *svc) ListProducts(ctx context.Context) error {
	return nil
}
