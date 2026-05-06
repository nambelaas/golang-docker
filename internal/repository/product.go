package repository

import (
	"database/sql"
	"fmt"

	"github.com/nambelaas/golang-docker/internal/model"
)

// ProductRepository handles all product-related database operations
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll retrieves all products
func (r *ProductRepository) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	var p model.Product

	err := r.db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product: %w", err)
	}

	return &p, nil
}

// Create inserts a new product
func (r *ProductRepository) Create(p *model.Product) error {
	err := r.db.QueryRow(
		"INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
		p.Name, p.Price,
	).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// Update updates an existing product
func (r *ProductRepository) Update(id int, p *model.Product) error {
	result, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2 WHERE id = $3",
		p.Name, p.Price, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Delete removes a product
func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
