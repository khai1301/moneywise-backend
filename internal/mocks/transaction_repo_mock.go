package mocks

import (
	"time"
	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(tx *models.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(id, userID string) (*models.Transaction, error) {
	args := m.Called(id, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) FindByUserID(userID string, startDate, endDate time.Time, limit, offset int) ([]models.Transaction, int64, error) {
	args := m.Called(userID, startDate, endDate, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Transaction), int64(args.Int(1)), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *MockTransactionRepository) Update(tx *models.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(id, userID string) error {
	args := m.Called(id, userID)
	return args.Error(0)
}
