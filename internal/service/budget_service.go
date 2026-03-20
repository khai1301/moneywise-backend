package service

import (
	"errors"
	"time"

	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/khai1301/moneywise-backend/internal/repository"
	"gorm.io/gorm"
)

// BudgetWithStats extends Budget with real spending stats for the requested month.
type BudgetWithStats struct {
	models.Budget
	Spent      float64 `json:"spent"`
	Remaining  float64 `json:"remaining"`
	Percentage float64 `json:"percentage"` // 0-100+
}

type BudgetService interface {
	CreateBudget(userID, categoryID, month, note string, amount float64) (*models.Budget, error)
	GetBudgetsWithStats(userID, month string, db *gorm.DB) ([]BudgetWithStats, error)
	UpdateBudget(id, userID, categoryID, month, note string, amount float64) (*models.Budget, error)
	DeleteBudget(id, userID string) error
}

type budgetService struct {
	budgetRepo   repository.BudgetRepository
	categoryRepo repository.CategoryRepository
	db           *gorm.DB
}

func NewBudgetService(budgetRepo repository.BudgetRepository, categoryRepo repository.CategoryRepository, db *gorm.DB) BudgetService {
	return &budgetService{budgetRepo: budgetRepo, categoryRepo: categoryRepo, db: db}
}

func (s *budgetService) CreateBudget(userID, categoryID, month, note string, amount float64) (*models.Budget, error) {
	if amount <= 0 {
		return nil, errors.New("số tiền ngân sách phải lớn hơn 0")
	}
	category, err := s.categoryRepo.FindByID(categoryID, userID)
	if err != nil || category == nil {
		return nil, errors.New("danh mục không hợp lệ")
	}

	budget := &models.Budget{
		UserID:     userID,
		CategoryID: categoryID,
		Amount:     amount,
		Month:      month,
		Note:       note,
	}
	if err := s.budgetRepo.Create(budget); err != nil {
		return nil, err
	}
	return s.budgetRepo.FindByID(budget.ID, userID)
}

func (s *budgetService) GetBudgetsWithStats(userID, month string, db *gorm.DB) ([]BudgetWithStats, error) {
	budgets, err := s.budgetRepo.FindByUserAndMonth(userID, month)
	if err != nil {
		return nil, err
	}

	// Parse month boundaries for the spending query
	startOfMonth, _ := time.Parse("2006-01", month)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	result := make([]BudgetWithStats, 0, len(budgets))
	for _, b := range budgets {
		var spent float64
		db.Model(&models.Transaction{}).
			Where("user_id = ? AND category_id = ? AND type = 'expense' AND date >= ? AND date <= ?",
				userID, b.CategoryID, startOfMonth, endOfMonth).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&spent)

		remaining := b.Amount - spent
		pct := 0.0
		if b.Amount > 0 {
			pct = (spent / b.Amount) * 100
		}

		result = append(result, BudgetWithStats{
			Budget:     b,
			Spent:      spent,
			Remaining:  remaining,
			Percentage: pct,
		})
	}
	return result, nil
}

func (s *budgetService) UpdateBudget(id, userID, categoryID, month, note string, amount float64) (*models.Budget, error) {
	budget, err := s.budgetRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if budget == nil {
		return nil, errors.New("không tìm thấy ngân sách")
	}
	if amount > 0 {
		budget.Amount = amount
	}
	if categoryID != "" {
		budget.CategoryID = categoryID
	}
	if month != "" {
		budget.Month = month
	}
	budget.Note = note

	if err := s.budgetRepo.Update(budget); err != nil {
		return nil, err
	}
	return s.budgetRepo.FindByID(id, userID)
}

func (s *budgetService) DeleteBudget(id, userID string) error {
	return s.budgetRepo.Delete(id, userID)
}
