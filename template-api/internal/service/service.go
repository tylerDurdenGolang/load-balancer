package service

import (
	"context"
	"template-api/internal/models"
	"template-api/internal/repository"
)

type ProductService struct {
	repo repository.IRepository
}

func NewProductService(repo repository.IRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateItem(
	ctx context.Context,
	name string,
	description string,
	price float64,
	stock int64,
) (int64, error) {
	return s.repo.CreateProduct(
		ctx,
		name,
		description,
		price,
		stock,
	)
}

func (s *ProductService) GetAllItems(ctx context.Context, limit int64) ([]models.Item, error) {
	return s.repo.GetAllProducts(ctx, limit)
}

func (s *ProductService) GetItemById(ctx context.Context, id int64) (*models.Item, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) UpdateItem(
	ctx context.Context,
	id int64,
	name *string,
	description *string,
	price *float64,
	stock *int64,
) error {
	panic("implement me")
}

func (s *ProductService) DeleteItem(ctx context.Context, id int64) error {
	panic("implement me")
}
