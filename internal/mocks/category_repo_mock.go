package mocks

import (
	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByUserID(userID string) ([]models.Category, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepository) FindByID(id string, userID string) (*models.Category, error) {
	args := m.Called(id, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepository) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id string, userID string) error {
	args := m.Called(id, userID)
	return args.Error(0)
}
