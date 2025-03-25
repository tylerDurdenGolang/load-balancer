package service

import (
	"context"
	"template-api/internal/models"
)

type IService interface {
	CreateItem(
		ctx context.Context,
		name string,
		description string,
		price float64,
		stock int64,
	) (int64, error)
	GetAllItems(ctx context.Context, limit int64) ([]models.Item, error)
	GetItemById(ctx context.Context, id int64) (models.Item, error)
	UpdateItem(ctx context.Context, item *models.UpdateItem) error
	DeleteItem(ctx context.Context, id int64) error
}
