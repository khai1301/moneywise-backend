package service

import (
	"errors"
	"time"

	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/khai1301/moneywise-backend/internal/repository"
)

type TransactionService interface {
	CreateTransaction(userID, categoryID, title string, amount float64, txType string, date time.Time, paymentMethod, note string) (*models.Transaction, error)
	GetTransactionsByUser(userID string, startDate, endDate time.Time, categoryID, txType string, limit, offset int) ([]models.Transaction, int64, error)
	UpdateTransaction(id, userID, categoryID, title string, amount float64, txType string, date time.Time, paymentMethod, note string) (*models.Transaction, error)
	DeleteTransaction(id, userID string) error
	GetTransactionByID(id, userID string) (*models.Transaction, error)
}

type transactionService struct {
	txRepo       repository.TransactionRepository
	categoryRepo repository.CategoryRepository
}

func NewTransactionService(txRepo repository.TransactionRepository, categoryRepo repository.CategoryRepository) TransactionService {
	return &transactionService{
		txRepo:       txRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *transactionService) CreateTransaction(userID, categoryID, title string, amount float64, txType string, date time.Time, paymentMethod, note string) (*models.Transaction, error) {
	if title == "" || txType == "" {
		return nil, errors.New("Tiêu đề và phân loại(type) không được rỗng")
	}

	// Lấy category để kiểm tra quyền sở hữu
	category, err := s.categoryRepo.FindByID(categoryID, userID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("Danh mục không hợp lệ hoặc bạn không có quyền sử dụng")
	}

	tx := &models.Transaction{
		UserID:        userID,
		CategoryID:    categoryID,
		Title:         title,
		Amount:        amount,
		Type:          txType,
		Date:          date,
		PaymentMethod: paymentMethod,
		Note:          note,
	}

	if err := s.txRepo.Create(tx); err != nil {
		return nil, err
	}

	// Đọc lại để có được thông tin Category Preload đầy đủ
	return s.txRepo.FindByID(tx.ID, userID)
}

func (s *transactionService) GetTransactionsByUser(userID string, startDate, endDate time.Time, categoryID, txType string, limit, offset int) ([]models.Transaction, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.txRepo.FindByUserID(userID, startDate, endDate, categoryID, txType, limit, offset)
}

func (s *transactionService) GetTransactionByID(id, userID string) (*models.Transaction, error) {
	return s.txRepo.FindByID(id, userID)
}

func (s *transactionService) UpdateTransaction(id, userID, categoryID, title string, amount float64, txType string, date time.Time, paymentMethod, note string) (*models.Transaction, error) {
	tx, err := s.txRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("Không tìm thấy Giao dịch")
	}

	// Nếu đổi category thì phải check category mới
	if tx.CategoryID != categoryID && categoryID != "" {
		category, err := s.categoryRepo.FindByID(categoryID, userID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("Danh mục mới không được phép sử dụng")
		}
		tx.CategoryID = categoryID
	}

	if title != "" {
		tx.Title = title
	}
	if txType != "" {
		tx.Type = txType
	}
	tx.Amount = amount
	
	if !date.IsZero() {
		tx.Date = date
	}
	tx.PaymentMethod = paymentMethod
	tx.Note = note

	if err := s.txRepo.Update(tx); err != nil {
		return nil, err
	}

	return s.txRepo.FindByID(tx.ID, userID)
}

func (s *transactionService) DeleteTransaction(id, userID string) error {
	return s.txRepo.Delete(id, userID)
}
