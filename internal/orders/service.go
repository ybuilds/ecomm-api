package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/ybuilds/ecomm-api/internal/adapters/postgresql/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product doesn't have enough stock")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (svc *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, errors.New("customer ID is required")
	}

	if len(tempOrder.Items) == 0 {
		return repo.Order{}, errors.New("at least one item is required")
	}

	//transaction creation
	tx, err := svc.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)
	qtx := svc.repo.WithTx(tx)

	//create order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range tempOrder.Items {
		product, err := svc.repo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		if err != nil {
			return repo.Order{}, err
		}

		//update item stock
	}

	err = tx.Commit(ctx)
	if err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
