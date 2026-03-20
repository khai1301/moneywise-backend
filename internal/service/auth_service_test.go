package service

import (
	"testing"
	"github.com/spf13/viper"

	"github.com/khai1301/moneywise-backend/internal/mocks"
	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/khai1301/moneywise-backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	viper.Set("JWT_SECRET", "test-secret")
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := NewAuthService(mockRepo)

	existingUser := &models.User{Email: "test@example.com"}

	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	user, err := authService.Register("Test", "test@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, "email đã được sử dụng", err.Error())
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := NewAuthService(mockRepo)

	mockRepo.On("FindByEmail", "new@example.com").Return(nil, nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := authService.Register("Test", "new@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "Test", user.Name)

	assert.NotEqual(t, "password123", user.PasswordHash)
    
	isValid := utils.CheckPasswordHash("password123", user.PasswordHash)
	assert.True(t, isValid)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := NewAuthService(mockRepo)

	hashedPassword, _ := utils.HashPassword("correct_password")
	existingUser := &models.User{
		ID:           "test-id",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	token, err := authService.Login("test@example.com", "wrong_password")

	assert.Error(t, err)
	assert.Equal(t, "email hoặc mật khẩu không chính xác", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := NewAuthService(mockRepo)

	hashedPassword, _ := utils.HashPassword("correct_password")
	existingUser := &models.User{
		ID:           "test-id",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	token, err := authService.Login("test@example.com", "correct_password")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}
