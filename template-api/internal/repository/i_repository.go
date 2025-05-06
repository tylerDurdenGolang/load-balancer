package repository

import (
	"context"
	"template-api/internal/models"
)

type IRepository interface {
	CreateProduct(
		ctx context.Context,
		name string,
		description string,
		price float64,
		stock int64,
	) (int64, error)
	GetProductByID(ctx context.Context, id int64) (*models.Item, error)
	GetAllProducts(ctx context.Context, limit int64) ([]models.Item, error)
}
