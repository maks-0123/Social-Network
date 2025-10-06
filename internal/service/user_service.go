package service

import (
	"context"
	"errors"
	"my-social-network/internal/models"
	"my-social-network/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
func (s *UserService) RegisterUser(email, username, password string) error {
	// Проверяем, занят ли email - используем репозиторий
	if len(password) < 6 {
		return errors.New("пароль должен быть не меньше 6 цифр")
	}
	if len(username) < 3 {
		return errors.New("имя пользователя должно быть не менее 3 символов")
	}
	ctx := context.Background()
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)

	if existingUser != nil {
		return errors.New("email уже занят")
	}

	existingUserByUsername, _ := s.userRepo.FindByUsername(ctx, username)
	if existingUserByUsername != nil {
		return errors.New("имя пользователя занято")
	}
	user := &models.User{
		Email:    email,
		Username: username,
		Password: password,
	}
	return s.userRepo.CreateUser(ctx, user)
}
func (s *UserService) LoginUser(email, password string) (*models.User, error) {
	ctx := context.Background()

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("певерный пароль")
	}

	user.Password = ""
	return user, nil
}
