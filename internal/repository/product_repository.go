package repository

import (
	"context"
	"database/sql"

	"github.com/prawirdani/go-midtrans-example/internal/entity"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]entity.Product, error)
	GetProductByID(ctx context.Context, id int) (entity.Product, error)
	SaveChanges(ctx context.Context, p entity.Product) error
}

type productRepository struct {
	dbConn *sql.DB
}

func NewProductRepository(dbConn *sql.DB) ProductRepository {
	return &productRepository{
		dbConn: dbConn,
	}
}

func (r *productRepository) GetProducts(ctx context.Context) ([]entity.Product, error) {
	products := make([]entity.Product, 0)

	rows, err := r.dbConn.QueryContext(ctx, selectProductQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, id int) (entity.Product, error) {
	var product entity.Product
	query := selectProductQuery + " WHERE id = ?"
	err := r.dbConn.QueryRowContext(ctx, query, id).
		Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Product{}, entity.ErrProductNotFound
		}
		return entity.Product{}, err
	}

	return product, nil
}

func (r *productRepository) SaveChanges(ctx context.Context, p entity.Product) error {
	query := "UPDATE products SET name=?, price=?  WHERE id = ?"
	_, err := r.dbConn.ExecContext(ctx, query, p.Name, p.Price, p.ID)
	return err
}

const selectProductQuery = "SELECT id, name, price FROM products"
