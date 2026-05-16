package services

import (
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/google/uuid"
	"github.com/strconvitoa/martian-service/internal/core/domain"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user domain.User) (domain.User, error) {
	id := uuid.New().String()
	user.ID = id
	_, err := s.repo.Save(user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
func (s *userService) GetByEmail(email string) (domain.User, error) {

	savedUser, err := s.repo.Get(email)
	if err != nil {
		return domain.User{}, err
	}
	return savedUser, nil
}
func (s *userService) ChangeStatus(user domain.User) (domain.User, error) {

	savedUser, err := s.repo.UpdateStatus(user)
	if err != nil {
		return domain.User{}, err
	}
	return savedUser, nil
}
func (s *userService) ChangePassword(user domain.User) (domain.User, error) {

	savedUser, err := s.repo.UpdatePassword(user)
	if err != nil {
		return domain.User{}, err
	}
	return savedUser, nil
}
