package handler

import (
	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type OrgHandler struct {
	orgsvc ports.OrgService
	usrsvc ports.UserService
}

func NewOrgHandler(orgsvc ports.OrgService, usrsvc ports.UserService) *OrgHandler {
	return &OrgHandler{orgsvc: orgsvc, usrsvc: usrsvc}
}

func (h *OrgHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		EntName   string `json:"entname"`
		TradeName string `json:"trade_name"`
		Password  string `json:"password"`
		Status    string `json:"status"`
		Phone     string `json:"phone"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Bad Request")
	}
	ent := domain.Org{Name: req.EntName, TradeName: req.TradeName}
	savedOrg, err := h.orgsvc.CreateOrg(ent)
	usr := domain.User{Name: req.Name, Email: req.Email, Password: req.Password, OrgID: savedOrg.ID, Role: "admin", Status: "active", Phone: req.Phone}
	saveduser, err := h.usrsvc.CreateUser(usr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}
	resp := saveduser.ToResponse()
	return c.JSON(fiber.Map{"user": resp})
}
