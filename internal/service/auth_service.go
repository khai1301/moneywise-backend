package service

import (
	"errors"

	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/khai1301/moneywise-backend/internal/repository"
	"github.com/khai1301/moneywise-backend/pkg/utils"
)

type AuthService interface {
	Register(name, email, password string) (*models.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(name, email, password string) (*models.User, error) {
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email đã được sử dụng")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("email hoặc mật khẩu không chính xác")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("email hoặc mật khẩu không chính xác")
	}

	token, err := utils.GenerateJWT(user.ID) 
	if err != nil {
		return "", err
	}

	return token, nil
}