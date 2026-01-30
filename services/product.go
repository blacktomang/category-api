package services

import (
	"category-api/models"
	"category-api/repositories"
	"errors"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	return s.repo.GetProducts()
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	return s.repo.GetProductByID(id)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	if product.Name == "" || product.Price == 0 || product.Stock == 0 || product.CategoryID == 0 {
		return errors.New("Invalid Payload")
	}

	res, _ := s.repo.GetProductByName(product.Name)

	if res != nil {
		return errors.New("Product already exists")
	}

	return s.repo.CreateProduct(product)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	if product.Name == "" && product.Price == 0 && product.Stock == 0 && product.CategoryID == 0 {
		return errors.New("Invalid Payload")
	}
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	if id <= 0 {
		return errors.New("Invalid ID")
	}
	return s.repo.DeleteProduct(id)
}
