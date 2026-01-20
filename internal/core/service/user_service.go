package service

import (
	"fmt"

	"go-test/go-test/internal/core/domain"
	"go-test/go-test/internal/core/repository"
)

type UserService interface {
	GetUserByID(id int) (*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetUserByID(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", id)
	}

	realId := 2 * id

	user, err := s.repo.GetByID(realId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) CreateUser(user *domain.User) (*domain.User, error) {
	if user.Name == "" {
		return nil, fmt.Errorf("user name is required")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email is required")
	}

	if user.Name == "Palinho" {
		user.Username = "usuario PCD"
	}

	createdUser, err := s.repo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}
