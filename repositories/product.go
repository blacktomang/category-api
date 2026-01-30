package repositories

import (
	"category-api/models"
	"database/sql"
	"errors"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"

	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetProducts() ([]models.Product, error) {
	products := make([]models.Product, 0)
	query := "SELECT id, name, price, stock, category_id FROM products LIMIT 10"

	rows, query_err := r.db.Query(query)

	if query_err != nil {
		return nil, query_err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(id int) (*models.Product, error) {
	var product models.Product
	query := "SELECT id, name, price, stock, category_id FROM products WHERE id = $1"

	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)

	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetProductByName(name string) (*models.Product, error) {
	var product models.Product
	query := "SELECT id FROM products WHERE name = $1"

	err := r.db.QueryRow(query, name).Scan(&product.ID)

	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"

	result, err := r.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Product not found")
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Product not found")
	}
	return nil
}
