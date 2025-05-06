package repository

import (
	"context"
	"template-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(
	ctx context.Context,
	name string,
	description string,
	price float64,
	stock int64,
) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, createProductQuery, name, description, price, stock).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int64) (*models.Item, error) {
	product := &models.Item{}
	err := r.db.QueryRow(ctx, getProductByIdQuery, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	return product, err
}

func (r *ProductRepository) GetAllProducts(ctx context.Context, limit int64) ([]models.Item, error) {
	rows, err := r.db.Query(ctx, getAllProductsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Item
	for rows.Next() {
		var product models.Item
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
