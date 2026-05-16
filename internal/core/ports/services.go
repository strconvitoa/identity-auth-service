package ports

import "github.com/strconvitoa/martian-service/internal/core/domain"

type AuthService interface {
	Validate(auth domain.Auth) (bool, error)
	Authenticate(auth domain.Auth) (domain.Auth, error)
	Reset(auth domain.Auth, timeToExpire string) (domain.Auth, error)
	Remove(id string) (bool, error)
}

type UserService interface {
	CreateUser(user domain.User) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	ChangePassword(user domain.User) (domain.User, error)
	ChangeStatus(user domain.User) (domain.User, error)
}

type OrgService interface {
	CreateOrg(Org domain.Org) (domain.Org, error)
}
type IntakeService interface {
	CreateIntake(intake domain.Intake) (domain.Intake, error)
}
