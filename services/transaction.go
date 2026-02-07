package services

import (
	"category-api/models"
	"category-api/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) ReportToday() (*models.Report, error) {
	result, err := s.repo.Report([]string{time.Now().Format("2006-01-02")})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TransactionService) ReportRange(dates []string) (*models.Report, error) {
	result, err := s.repo.Report(dates)
	if err != nil {
		return nil, err
	}
	return result, nil
}
