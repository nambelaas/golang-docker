package model

// Product represents a product entity
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required"`
}

// CreateProductRequest is the request body for creating a product
type CreateProductRequest struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required,gt=0"`
}

// UpdateProductRequest is the request body for updating a product
type UpdateProductRequest struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required,gt=0"`
}
