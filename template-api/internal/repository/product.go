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

func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Item) error {
	query := `INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.Stock).Scan(&product.ID)
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (*models.Item, error) {
	product := &models.Item{}
	query := `SELECT id, name, description, price, stock FROM products WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	return product, err
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Item, error) {
	query := `SELECT id, name, description, price, stock FROM products`
	rows, err := r.db.Query(ctx, query)
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
