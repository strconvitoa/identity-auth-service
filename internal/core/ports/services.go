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
	FindByEmail(email string) (domain.User, error)
	ChangePassword(user domain.User) (domain.User, error)
	ChangeStatus(user domain.User) (domain.User, error)
	IsExistingUser(email string) (bool, error)
	RemoveByEmail(email string) error
	UserExists(email string, password string) (bool, error)
	FindPasswordByEmail(email string) (string, error)
	FindAllByOrgID(orgID string) ([]domain.User, error)
}

type OrgService interface {
	CreateOrg(org domain.Org) (domain.Org, error)
	OrgExists(org domain.Org) (bool, error)
	OrgExistsByID(id string) (bool, error)
}
type LeadService interface {
	CreateLead(Lead domain.Lead) (domain.Lead, error)
	FindLeadByStatus(org_id string, status string) ([]domain.Lead, error)
}

type EmailService interface {
	SendEmail(toEmail string, subject string, body string) error
	GetVerificationEmailBody(otpCode string) string
	GetWelcomeEmailBody(otpCode string) string
}
