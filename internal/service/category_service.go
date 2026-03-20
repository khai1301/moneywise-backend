package service

import (
	"errors"

	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/khai1301/moneywise-backend/internal/repository"
)

type CategoryService interface {
	CreateCategory(userID, name, catType, icon, color, description string) (*models.Category, error)
	GetCategoriesByUser(userID string) ([]models.Category, error)
	UpdateCategory(id, userID, name, catType, icon, color, description string) (*models.Category, error)
	DeleteCategory(id, userID string) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(userID, name, catType, icon, color, description string) (*models.Category, error) {
	if name == "" || catType == "" {
		return nil, errors.New("Tên và Loại (income/expense) không được để trống")
	}

	category := &models.Category{
		UserID:      userID,
		Name:        name,
		Type:        catType,
		Icon:        icon,
		Color:       color,
		Description: description,
		IsSystem:    false,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategoriesByUser(userID string) ([]models.Category, error) {
	return s.categoryRepo.FindByUserID(userID)
}

func (s *categoryService) UpdateCategory(id, userID, name, catType, icon, color, description string) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("Không tìm thấy Danh mục")
	}
	if category.IsSystem {
		return nil, errors.New("Không thể chỉnh sửa danh mục mặc định của hệ thống")
	}

	category.Name = name
	category.Type = catType
	category.Icon = icon
	category.Color = color
	category.Description = description

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id, userID string) error {
	return s.categoryRepo.Delete(id, userID)
}
