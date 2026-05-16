package handler

import (
	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc     ports.UserService
	authsvc ports.AuthService
}

func NewUserHandler(svc ports.UserService, authsvc ports.AuthService) *UserHandler {
	return &UserHandler{svc: svc, authsvc: authsvc}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name     string          `json:"name"`
		Email    string          `json:"email"`
		Password string          `json:"password"`
		OrgID    string          `json:"org_id"`
		Role     domain.Roles    `json:"role"`
		Status   domain.Statuses `json:"status"`
		Phone    string          `json:"phone"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Bad Request")
	}
	usr := domain.User{Name: req.Name, Email: req.Email, Password: req.Password, OrgID: req.OrgID, Role: req.Role, Status: "pending", Phone: req.Phone}
	savedusr, err := h.svc.CreateUser(usr)
	ent := domain.Auth{Email: savedusr.Email, UserID: savedusr.ID}
	_, _ = h.authsvc.Reset(ent, "24h")
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}
	resp := savedusr.ToResponse()

	return c.JSON(fiber.Map{"user": resp})
}
