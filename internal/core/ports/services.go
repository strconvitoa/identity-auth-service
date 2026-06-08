package ports

import "github.com/strconvitoa/martian-service/internal/core/domain"

type AuthService interface {
	Validate(auth domain.Auth) (bool, error)
	Authenticate(auth domain.Auth) (domain.Auth, error)
	Reset(auth domain.Auth, timeToExpire string) (domain.Auth, error)
	Remove(id string) (bool, error)
	RemoveByEmail(email string) error
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type UserService interface {
	CreateUser(user domain.User) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	ChangePassword(user domain.User) (domain.User, error)
	ChangeStatus(user domain.User) (domain.User, error)
	IsExistingUser(email string) (bool, error)
	RemoveByEmail(email string) error
}

type OrgService interface {
	CreateOrg(org domain.Org) (domain.Org, error)
	OrgExists(org domain.Org) (bool, error)
	OrgExistsByID(id string) (bool, error)
}
type IntakeService interface {
	CreateIntake(intake domain.Intake) (domain.Intake, error)
}

type EmailService interface {
	SendEmail(toEmail string, subject string, body string) error
	GetVerificationEmailBody(otpCode string) string
	GetWelcomeEmailBody(otpCode string) string
}
