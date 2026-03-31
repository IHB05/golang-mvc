package repositories

import (
	"belajar-crud-mvc/models"

	"gorm.io/gorm"
)

// UserRepository mendefinisikan semua operasi database untuk user
type UserRepository interface {
	FindAll(page, limit int, search string) ([]models.User, int64, error)
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User, updates map[string]interface{}) error
	Delete(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// FindAll mengambil semua user dengan pagination dan search
func (r *userRepository) FindAll(page, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.User{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindByID mencari user berdasarkan ID
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email (untuk cek duplikat)
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create menyimpan user baru ke database
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update mengupdate field tertentu dari user
func (r *userRepository) Update(user *models.User, updates map[string]interface{}) error {
	return r.db.Model(user).Updates(updates).Error
}

// Delete soft delete user
func (r *userRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}
