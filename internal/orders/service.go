package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/luizflm/go-ecom-api/internal/adapters/postgresql/sqlc"
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

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range tempOrder.Items {
		_, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				return repo.Order{}, ErrProductNotFound
			default:
				return repo.Order{}, err
			}
		}

		product, err := qtx.ReduceProductStock(ctx, repo.ReduceProductStockParams{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				return repo.Order{}, ErrProductNoStock
			default:
				return repo.Order{}, err
			}
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      order.ID,
			ProductID:    product.ID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
