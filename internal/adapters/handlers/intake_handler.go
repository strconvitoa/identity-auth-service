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
		Name    string `json:"name"`
		Email   string `json:"email"`
		Issue   string `json:"issue"`
		Phone   string `json:"phone"`
		OrgID   string `json:"org_id"`
		Summary string `json:"summary,omitempty"`
		Lang    string `json:"lang,omitempty"`
	}
	res := domain.Envelope[domain.Intake]{Success: false, Data: domain.Intake{}, Error: ""}
	if err := c.BodyParser(&req); err != nil {
		res.Error = "Bad Request"
		return c.Status(400).JSON(res)
	}
	usr := domain.Intake{Name: req.Name, Email: req.Email, Phone: req.Phone, Issue: req.Issue, OrgID: req.OrgID, Summary: req.Summary, Lang: req.Lang}
	intake, err := h.svc.CreateIntake(usr)
	if err != nil {
		res.Success = false
		res.Data = intake
		res.Error = "Error creating intake"
		return c.Status(401).JSON(res)
	}
	res.Success = true
	res.Data = intake
	res.Error = ""
	return c.Status(200).JSON(res)
}
