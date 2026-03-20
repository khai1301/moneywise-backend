package repository

import (
	"errors"
	"time"

	"github.com/khai1301/moneywise-backend/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	FindByID(id, userID string) (*models.Transaction, error)
	FindByUserID(userID string, startDate, endDate time.Time, categoryID, txType string, limit, offset int) ([]models.Transaction, int64, error)
	Update(tx *models.Transaction) error
	Delete(id, userID string) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) FindByID(id, userID string) (*models.Transaction, error) {
	var tx models.Transaction
	result := r.db.Preload("Category").Where("id = ? AND user_id = ?", id, userID).First(&tx)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) FindByUserID(userID string, startDate, endDate time.Time, categoryID, txType string, limit, offset int) ([]models.Transaction, int64, error) {
	var txs []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("user_id = ?", userID)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if txType != "" && txType != "all" {
		query = query.Where("type = ?", txType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Order("date DESC").Limit(limit).Offset(offset).Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

func (r *transactionRepository) Update(tx *models.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *transactionRepository) Delete(id, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Không thể xóa giao dịch (không tồn tại hoặc không đủ quyền)")
	}
	return nil
}
