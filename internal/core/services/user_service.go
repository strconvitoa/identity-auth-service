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
	savusr, err := s.repo.Save(user)
	if err != nil {
		return domain.User{}, err
	}
	return savusr, nil
}
func (s *userService) IsExistingUser(email string) (bool, error) {

	exists, err := s.repo.Exists(email)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (s *userService) FindByEmail(email string) (domain.User, error) {

	savedUser, err := s.repo.Get(email)
	if err != nil {
		return domain.User{}, err
	}
	return savedUser, nil
}
func (s *userService) FindAllByOrgID(orgID string) ([]domain.User, error) {

	allusr, err := s.repo.SelectAllByOrgID(orgID)
	if err != nil {
		return []domain.User{}, err
	}
	return allusr, nil
}
func (s *userService) FindPasswordByEmail(email string) (string, error) {
	pword, err := s.repo.SelectPassword(email)
	if err != nil {
		return "", err
	}
	return pword, nil
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
func (s *userService) RemoveByEmail(email string) error {
	err := s.repo.DeleteByEmail(email)
	if err != nil {
		return err
	}
	return nil
}
func (s *userService) UserExists(email string, password string) (bool, error) {

	exists, err := s.repo.CredentialsMatch(email, password)
	if err != nil {
		return false, err
	}
	return exists, nil
}
