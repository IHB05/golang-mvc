package services

import (
	"belajar-crud-mvc/models"
	"belajar-crud-mvc/repositories"
	"errors"

	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProducts(page, limit int, search, category string) ([]models.Product, int64, int, error)
	GetProductByID(id uint) (*models.Product, error)
	CreateProduct(input models.CreateProductInput) (*models.Product, error)
	UpdateProduct(id uint, input models.UpdateProductInput) (*models.Product, error)
	DeleteProduct(id uint) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(page, limit int, search, category string) ([]models.Product, int64, int, error) {
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }

	products, total, err := s.repo.FindAll(page, limit, search, category)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return products, total, totalPages, nil
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("failed to fetch product")
	}
	return product, nil
}

func (s *productService) CreateProduct(input models.CreateProductInput) (*models.Product, error) {
	product := &models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		Category:    input.Category,
	}
	if err := s.repo.Create(product); err != nil {
		return nil, errors.New("failed to create product")
	}
	return product, nil
}

func (s *productService) UpdateProduct(id uint, input models.UpdateProductInput) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("failed to fetch product")
	}

	updates := map[string]interface{}{}
	if input.Name != nil        { updates["name"] = *input.Name }
	if input.Description != nil { updates["description"] = *input.Description }
	if input.Price != nil       { updates["price"] = *input.Price }
	if input.Stock != nil       { updates["stock"] = *input.Stock }
	if input.Category != nil    { updates["category"] = *input.Category }

	if len(updates) == 0 {
		return nil, errors.New("no fields provided to update")
	}

	if err := s.repo.Update(product, updates); err != nil {
		return nil, errors.New("failed to update product")
	}

	updated, _ := s.repo.FindByID(id)
	return updated, nil
}

func (s *productService) DeleteProduct(id uint) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return errors.New("failed to fetch product")
	}
	return s.repo.Delete(product)
}
