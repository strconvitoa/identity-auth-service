package handler

import (
	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type IntakeHandler struct {
	svc ports.IntakeService
}

func NewIntakeHandler(svc ports.IntakeService) *IntakeHandler {
	return &IntakeHandler{svc: svc}
}

func (h *IntakeHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Issue string `json:"issue"`
		Phone string `json:"phone"`
		OrgID string `json:"org_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Bad Request")
	}
	usr := domain.Intake{Name: req.Name, Email: req.Email, Phone: req.Phone, Issue: req.Issue, OrgID: req.OrgID}
	// NEED FIXING
	intake, err := h.svc.CreateIntake(usr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"intake": intake})
}
