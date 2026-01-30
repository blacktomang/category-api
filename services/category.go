package services

import (
	"category-api/models"
	"category-api/repositories"
	"errors"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo}
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	return s.repo.GetCategories()
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	return s.repo.GetCategoryByID(id)
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	if category.Name == "" || category.Description == "" {
		return errors.New("Invalid Payload")
	}

	result, _ := s.repo.GetCategoryByName(category.Name)

	if result != nil {
		return errors.New("Category already exists")
	}
	return s.repo.CreateCategory(category)
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	if category.Name == "" && category.Description == "" {
		return errors.New("Invalid Payload")
	}
	return s.repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	if id <= 0 {
		return errors.New("Invalid ID")
	}
	return s.repo.DeleteCategory(id)
}
