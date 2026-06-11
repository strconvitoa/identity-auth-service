package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/strconvitoa/martian-service/internal/core/ports"
	"golang.org/x/crypto/bcrypt"

	"github.com/strconvitoa/martian-service/internal/core/domain"
)

type authService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) ports.AuthService {
	return &authService{repo: repo}
}

func (s *authService) Validate(auth domain.Auth) (bool, error) {
	exists, err := s.repo.Exists(auth)
	if exists == false {
		return false, errors.New("token is expired")
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *authService) Authenticate(auth domain.Auth) (domain.Auth, error) {
	savedauth, err := s.repo.Get(auth)
	if err != nil {
		return domain.Auth{}, err
	}
	return savedauth, nil
}
func (s *authService) Reset(auth domain.Auth, timeToExpire string) (domain.Auth, error) {
	id := uuid.New().String()
	exp, err := s.generateExpiresAtTime(timeToExpire)
	token, _ := s.generateToken()
	auth.Token = token
	auth.ID = id
	auth.ExpiresAt = exp
	in, err := s.repo.Insert(auth)
	if err != nil {
		return domain.Auth{}, err
	}
	return in, nil
}
func (s *authService) Remove(id string) (bool, error) {
	_, err := s.repo.Delete(id)
	if err != nil {
		return false, err
	}
	return true, err
}
func (s *authService) RemoveByEmail(email string) error {
	err := s.repo.DeleteByEmail(email)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) generateExpiresAtTime(timeToExpire string) (string, error) {
	// 1. Parse the string into a duration (e.g., "15m" -> 15 minutes)
	duration, err := time.ParseDuration(timeToExpire)
	if err != nil {
		return "", fmt.Errorf("invalid duration format: %w", err)
	}

	// 2. Add the duration to the current time
	expirationTime := time.Now().Add(duration)

	// 3. Format as a string (RFC3339 is standard for APIs and DBs)
	// Example output: "2024-05-14T15:04:05Z"
	return expirationTime.Format(time.RFC3339), nil
}

func (s *authService) generateToken() (string, error) {
	// Max value is 1000000, so Int returns a number from 0 to 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format as a 6-digit string, padding with leading zeros if necessary (e.g., 004321)
	return fmt.Sprintf("%06d", n), nil
}

func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plain text password with a bcrypt hash
func (s *authService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
