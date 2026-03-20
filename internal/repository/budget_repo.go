package repository

import (
	"errors"

	"github.com/khai1301/moneywise-backend/internal/models"
	"gorm.io/gorm"
)

type BudgetRepository interface {
	Create(budget *models.Budget) error
	FindByID(id, userID string) (*models.Budget, error)
	FindByUserAndMonth(userID, month string) ([]models.Budget, error)
	Update(budget *models.Budget) error
	Delete(id, userID string) error
}

type budgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(budget *models.Budget) error {
	return r.db.Create(budget).Error
}

func (r *budgetRepository) FindByID(id, userID string) (*models.Budget, error) {
	var budget models.Budget
	result := r.db.Preload("Category").Where("id = ? AND user_id = ?", id, userID).First(&budget)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) FindByUserAndMonth(userID, month string) ([]models.Budget, error) {
	var budgets []models.Budget
	if err := r.db.Preload("Category").
		Where("user_id = ? AND month = ?", userID, month).
		Order("created_at DESC").
		Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *budgetRepository) Update(budget *models.Budget) error {
	return r.db.Save(budget).Error
}

func (r *budgetRepository) Delete(id, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Budget{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("không tìm thấy ngân sách hoặc không có quyền xóa")
	}
	return nil
}
