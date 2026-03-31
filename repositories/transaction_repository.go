package repositories

import (
	"belajar-crud-mvc/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindAll(page, limit int, userID uint) ([]models.Transaction, int64, error)
	FindByID(id uint) (*models.Transaction, error)
	Create(transaction *models.Transaction) error
	UpdateStatus(transaction *models.Transaction, status models.TransactionStatus) error
	Delete(transaction *models.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) FindAll(page, limit int, userID uint) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.Transaction{})

	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		// Omit("User", "Items") supaya GORM tidak mencoba INSERT User yang sudah ada
		// Kalau tidak di-Omit, GORM akan coba simpan User baru -> duplicate error
		if err := tx.Omit("User", "Items").Create(transaction).Error; err != nil {
			return err
		}

		// Simpan setiap item satu per satu
		// Omit("Product") supaya GORM tidak mencoba INSERT Product yang sudah ada
		for i := range transaction.Items {
			transaction.Items[i].TransactionID = transaction.ID
			if err := tx.Omit("Product").Create(&transaction.Items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *transactionRepository) UpdateStatus(transaction *models.Transaction, status models.TransactionStatus) error {
	return r.db.Model(transaction).Update("status", status).Error
}

func (r *transactionRepository) Delete(transaction *models.Transaction) error {
	return r.db.Delete(transaction).Error
}