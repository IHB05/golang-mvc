package repositories

import (
	"belajar-crud-mvc/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(page, limit int, search, category string) ([]models.Product, int64, error)
	FindByID(id uint) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product, updates map[string]interface{}) error
	Delete(product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FindAll(page, limit int, search, category string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&models.Product{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product, updates map[string]interface{}) error {
	return r.db.Model(product).Updates(updates).Error
}

func (r *productRepository) Delete(product *models.Product) error {
	return r.db.Delete(product).Error
}
