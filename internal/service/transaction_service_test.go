package service

import (
	"testing"
	"time"

	"github.com/khai1301/moneywise-backend/internal/mocks"
	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_CreateTransaction_MissingFields(t *testing.T) {
	mockTxRepo := new(mocks.MockTransactionRepository)
	mockCatRepo := new(mocks.MockCategoryRepository)
	svc := NewTransactionService(mockTxRepo, mockCatRepo)

	tx, err := svc.CreateTransaction("u1", "c1", "", 100, "income", time.Now(), "Cash", "Note")
	assert.Error(t, err)
	assert.Equal(t, "Tiêu đề và phân loại(type) không được rỗng", err.Error())
	assert.Nil(t, tx)
}

func TestTransactionService_CreateTransaction_InvalidCategory(t *testing.T) {
	mockTxRepo := new(mocks.MockTransactionRepository)
	mockCatRepo := new(mocks.MockCategoryRepository)
	svc := NewTransactionService(mockTxRepo, mockCatRepo)

	mockCatRepo.On("FindByID", "evil-cat", "u1").Return(nil, nil)

	tx, err := svc.CreateTransaction("u1", "evil-cat", "Lương", 1000, "income", time.Now(), "", "")
	
	assert.Error(t, err)
	assert.Equal(t, "Danh mục không hợp lệ hoặc bạn không có quyền sử dụng", err.Error())
	assert.Nil(t, tx)
	mockCatRepo.AssertExpectations(t)
}

func TestTransactionService_CreateTransaction_Success(t *testing.T) {
	mockTxRepo := new(mocks.MockTransactionRepository)
	mockCatRepo := new(mocks.MockCategoryRepository)
	svc := NewTransactionService(mockTxRepo, mockCatRepo)

	validCat := &models.Category{ID: "c1", UserID: "u1", Name: "Food"}
	mockCatRepo.On("FindByID", "c1", "u1").Return(validCat, nil)
	mockTxRepo.On("Create", mock.AnythingOfType("*models.Transaction")).Return(nil)
	
	createdTx := &models.Transaction{ID: "tx1", UserID: "u1"}
	mockTxRepo.On("FindByID", mock.AnythingOfType("string"), "u1").Return(createdTx, nil)

	tx, err := svc.CreateTransaction("u1", "c1", "Lương", 1000, "income", time.Now(), "", "")

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	mockCatRepo.AssertExpectations(t)
	mockTxRepo.AssertExpectations(t)
}
