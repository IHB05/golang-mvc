package services

import (
	"belajar-crud-mvc/models"
	"belajar-crud-mvc/repositories"
	"errors"

	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers(page, limit int, search string) ([]models.User, int64, int, error)
	GetUserByID(id uint) (*models.User, error)
	CreateUser(input models.CreateUserInput) (*models.User, error)
	UpdateUser(id uint, input models.UpdateUserInput) (*models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers(page, limit int, search string) ([]models.User, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := s.repo.FindAll(page, limit, search)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return users, total, totalPages, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to fetch user")
	}
	return user, nil
}

func (s *userService) CreateUser(input models.CreateUserInput) (*models.User, error) {
	// Cek apakah email sudah dipakai
	existing, _ := s.repo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password, // catatan: di production harus di-hash dulu
		Phone:    input.Phone,
		Address:  input.Address,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *userService) UpdateUser(id uint, input models.UpdateUserInput) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to fetch user")
	}

	// Cek email baru tidak bentrok dengan user lain
	if input.Email != nil {
		existing, _ := s.repo.FindByEmail(*input.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email already used by another user")
		}
	}

	updates := map[string]interface{}{}
	if input.Name != nil    { updates["name"] = *input.Name }
	if input.Email != nil   { updates["email"] = *input.Email }
	if input.Phone != nil   { updates["phone"] = *input.Phone }
	if input.Address != nil { updates["address"] = *input.Address }

	if len(updates) == 0 {
		return nil, errors.New("no fields provided to update")
	}

	if err := s.repo.Update(user, updates); err != nil {
		return nil, errors.New("failed to update user")
	}

	updated, _ := s.repo.FindByID(id)
	return updated, nil
}

func (s *userService) DeleteUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to fetch user")
	}
	return s.repo.Delete(user)
}
