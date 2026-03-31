package services

import (
	"belajar-crud-mvc/models"
	"belajar-crud-mvc/repositories"
	"errors"
	"fmt"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

type TransactionService interface {
	GetAllTransactions(page, limit int, userID uint) ([]models.Transaction, int64, int, error)
	GetTransactionByID(id uint) (*models.Transaction, error)
	CreateTransaction(input models.CreateTransactionInput) (*models.Transaction, error)
	UpdateTransactionStatus(id uint, input models.UpdateTransactionStatusInput) (*models.Transaction, error)
	DeleteTransaction(id uint) error
}

type TransactionServiceParams struct {
	dig.In

	TransactionRepo repositories.TransactionRepository
	ProductRepo     repositories.ProductRepository
	UserRepo        repositories.UserRepository
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	productRepo     repositories.ProductRepository
	userRepo        repositories.UserRepository
}

func NewTransactionService(p TransactionServiceParams) TransactionService {
	return &transactionService{
		transactionRepo: p.TransactionRepo,
		productRepo:     p.ProductRepo,
		userRepo:        p.UserRepo,
	}
}

func (s *transactionService) GetAllTransactions(page, limit int, userID uint) ([]models.Transaction, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	transactions, total, err := s.transactionRepo.FindAll(page, limit, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return transactions, total, totalPages, nil
}

func (s *transactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) CreateTransaction(input models.CreateTransactionInput) (*models.Transaction, error) {
	// 1. Pastikan user ada
	_, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 2. Proses setiap item
	var items []models.TransactionItem
	var totalPrice float64

	for _, itemInput := range input.Items {
		// Cek produk ada
		product, err := s.productRepo.FindByID(itemInput.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("product not found")
			}
			return nil, err
		}

		// Cek stok cukup
		if product.Stock < itemInput.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s (stock: %d, requested: %d)",
				product.Name, product.Stock, itemInput.Quantity)
		}

		// Hitung subtotal
		subtotal := product.Price * float64(itemInput.Quantity)
		totalPrice += subtotal

		// Kurangi stok
		newStock := product.Stock - itemInput.Quantity
		if err := s.productRepo.Update(product, map[string]interface{}{"stock": newStock}); err != nil {
			return nil, fmt.Errorf("failed to update stock: %v", err)
		}

		items = append(items, models.TransactionItem{
			ProductID: product.ID,
			Quantity:  itemInput.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
		})
	}

	// 3. Buat transaksi
	transaction := &models.Transaction{
		UserID:     input.UserID,
		TotalPrice: totalPrice,
		Status:     models.StatusPending,
		Notes:      input.Notes,
		Items:      items,
	}

	// 4. Simpan ke database — error asli dikembalikan langsung
	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %v", err)
	}

	// 5. Ambil transaksi lengkap dengan relasi
	created, err := s.transactionRepo.FindByID(transaction.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created transaction: %v", err)
	}

	return created, nil
}

func (s *transactionService) UpdateTransactionStatus(id uint, input models.UpdateTransactionStatusInput) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	if transaction.Status == models.StatusCancelled {
		return nil, errors.New("cancelled transaction cannot be updated")
	}

	if err := s.transactionRepo.UpdateStatus(transaction, input.Status); err != nil {
		return nil, fmt.Errorf("failed to update status: %v", err)
	}

	updated, _ := s.transactionRepo.FindByID(id)
	return updated, nil
}

func (s *transactionService) DeleteTransaction(id uint) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return err
	}
	return s.transactionRepo.Delete(transaction)
}
