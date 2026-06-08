package domain

type Statuses string
type Roles string

const (
	RoleAdmin   Roles = "admin"
	RoleManager Roles = "manager"
	RoleMember  Roles = "member"
)
const (
	StatusPending   Statuses = "pending"
	StatusActive    Statuses = "active"
	StatusSuspended Statuses = "suspended"
)

type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Email    string   `json:"email"`
	OrgID    string   `json:"org_id"`
	Role     Roles    `json:"role"`
	Status   Statuses `json:"status"`
	Phone    string   `json:"phone"`
}

func (u User) ToResponse() UserResponseDTO {
	return UserResponseDTO{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Role:   string(u.Role),
		Status: string(u.Status),
		OrgID:  u.OrgID,
		Phone:  u.Phone,
	}
}

type UserResponseDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
	OrgID  string `json:"org_id"`
	Phone  string `json:"phone"`
}
