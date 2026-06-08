package ports

import "github.com/strconvitoa/martian-service/internal/core/domain"

type UserRepository interface {
	Save(user domain.User) (domain.User, error)
	Get(email string) (domain.User, error)
	UpdatePassword(user domain.User) (domain.User, error)
	UpdateStatus(user domain.User) (domain.User, error)
	Exists(email string) (bool, error)
	DeleteByEmail(email string) error
	//Delete(id string) error
}
type OrgRepository interface {
	Save(user domain.Org) (domain.Org, error)
	Exists(org domain.Org) (bool, error)
	ExistsByID(id string) (bool, error)
	// GetByName(email string) (domain.Org, error)
	// GetByID(id string) (domain.Org, error)
	// Delete(id string) error
}

type IntakeRepository interface {
	Save(user domain.Intake) (domain.Intake, error)
	// GetByName(email string) (domain.Org, error)
	// GetByID(id string) (domain.Org, error)
	// Delete(id string) error
}

type AuthRepository interface {
	Insert(auth domain.Auth) (domain.Auth, error)
	Get(auth domain.Auth) (domain.Auth, error)
	Exists(auth domain.Auth) (bool, error)
	Delete(id string) (bool, error)
	DeleteByEmail(email string) error
}
