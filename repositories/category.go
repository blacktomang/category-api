package repositories

import (
	"category-api/models"
	"database/sql"
	"errors"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) GetCategories() ([]models.Category, error) {
	categories := make([]models.Category, 0)
	query := "SELECT id, name, description FROM categories LIMIT 10"

	rows, query_err := r.db.Query(query)

	if query_err != nil {
		return nil, query_err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	var category models.Category
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	err := r.db.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("Category not found")
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	query := "SELECT id FROM categories WHERE name = $1"

	err := r.db.QueryRow(query, name).Scan(&category.ID)
	if err == sql.ErrNoRows {
		return nil, errors.New("Category not found")
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"

	result, err := r.db.Exec(query, category.Name, category.Description, category.ID)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Category not found")
	}
	return nil
}

func (r *CategoryRepository) DeleteCategory(id int) error {
	query := "DELETE FROM categories WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Category not found")
	}

	return nil
}
