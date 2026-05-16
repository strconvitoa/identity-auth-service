package services

import (
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"github.com/strconvitoa/martian-service/internal/core/ports"

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
	auth.ID = id
	auth.ExpiresAt = exp
	in, err := s.repo.Insert(auth)
	if err != nil {
		return domain.Auth{}, err
	}
	return in, nil
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

func (s *authService) SendEmail(toEmail string, subject string, body string) error {
	// 1. Configuration - Use environment variables for these!
	from := "your-email@gmail.com"
	password := "your-app-password" // Not your login password!
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 2. Format the message
	// SMTP messages require a specific format: "Subject: ... \n\n Body"
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// 3. Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 4. Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
func (s *authService) Remove(id string) (bool, error) {
	_, err := s.repo.Delete(id)
	if err != nil {
		return false, err
	}
	return true, err
}
