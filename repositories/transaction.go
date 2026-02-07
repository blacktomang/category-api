package repositories

import (
	"category-api/models"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	if len(items) < 1 {
		return nil, fmt.Errorf("at least 1 item is required for checkout")
	}
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err = tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
		fmt.Printf("Appended detail, new length: %d\n", len(details))
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		var detailID int
		err = tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)

		fmt.Printf("is it error %v\n", err)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID

		fmt.Printf("returned id : %v\n", detailID)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) Report(dates []string) (*models.Report, error) {
	if len(dates) == 0 {
		return nil, fmt.Errorf("at least 1 date is required for report")
	}

	if len(dates) == 1 {
		revenueAndTotalquery := "SELECT COUNT(id) as total_transactions, COALESCE(SUM(total_amount), 0) as total_revenue FROM transactions WHERE created_at =  $1"
		transactionResult, transationErr := repo.db.Query(revenueAndTotalquery, dates[0])
		if transationErr != nil {
			return nil, transationErr
		}
		defer transactionResult.Close()

		var report models.Report
		for transactionResult.Next() {
			transationErr := transactionResult.Scan(&report.TotalTransactions, &report.TotalRevenue)
			if transationErr != nil {
				return nil, transationErr
			}
		}

		mostPopularProductQuery := `SELECT td.product_id, p.name as product_name, SUM(td.quantity) as total_sold
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at = $1
		GROUP BY td.product_id, p.name
		ORDER BY total_sold DESC
		LIMIT 1`

		mostPopularProductResult, mostPopularProductErr := repo.db.Query(mostPopularProductQuery, dates[0])

		if mostPopularProductErr != nil {
			return nil, mostPopularProductErr
		}

		defer mostPopularProductResult.Close()

		for mostPopularProductResult.Next() {
			mostPopularProductResult.Scan(&report.MostPopularProduct.ProductID, &report.MostPopularProduct.ProductName, &report.MostPopularProduct.TotalSold)
		}

		return &report, nil
	} else if len(dates) == 2 {
		revenueAndTotalquery := "SELECT COUNT(id) as total_transactions, COALESCE(SUM(total_amount), 0) as total_revenue FROM transactions WHERE created_at >= $1 AND created_at <= $2"
		transactionResult, transationErr := repo.db.Query(revenueAndTotalquery, dates[0], dates[1])
		if transationErr != nil {
			return nil, transationErr
		}
		defer transactionResult.Close()

		var report models.Report
		for transactionResult.Next() {
			transationErr := transactionResult.Scan(&report.TotalTransactions, &report.TotalRevenue)
			if transationErr != nil {
				return nil, transationErr
			}
		}

		mostPopularProductQuery := `SELECT td.product_id, p.name as product_name, SUM(td.quantity) as total_sold
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY td.product_id, p.name
		ORDER BY total_sold DESC
		LIMIT 1`

		mostPopularProductResult, mostPopularProductErr := repo.db.Query(mostPopularProductQuery, dates[0], dates[1])

		if mostPopularProductErr != nil {
			return nil, mostPopularProductErr
		}

		defer mostPopularProductResult.Close()

		for mostPopularProductResult.Next() {
			mostPopularProductResult.Scan(&report.MostPopularProduct.ProductID, &report.MostPopularProduct.ProductName, &report.MostPopularProduct.TotalSold)
		}

		return &report, nil
	}
	// Implementation for daily report can be added here
	return nil, nil
}
