package handler

import (
	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type OrgHandler struct {
	orgsvc  ports.OrgService
	usrsvc  ports.UserService
	authsvc ports.AuthService
}

func NewOrgHandler(orgsvc ports.OrgService, usrsvc ports.UserService, authsvc ports.AuthService) *OrgHandler {
	return &OrgHandler{orgsvc: orgsvc, usrsvc: usrsvc, authsvc: authsvc}
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
	res := domain.Envelope[domain.User]{Success: false, Data: domain.User{}, Error: "Bad Request"}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(res)
	}
	org := domain.Org{Name: req.EntName, TradeName: req.TradeName}
	orgExist, err := h.orgsvc.OrgExists(org)
	if orgExist == true {
		res.Error = "Organization already exists"
		return c.Status(400).JSON(res)
	}
	savedOrg, err := h.orgsvc.CreateOrg(org)
	pword, err := h.authsvc.HashPassword(req.Password)
	if err != nil {
		res.Error = "Error with hash"
		return c.Status(401).JSON(res)
	}
	usr := domain.User{Name: req.Name, Email: req.Email, Password: pword, OrgID: savedOrg.ID, Role: "admin", Status: "active", Phone: req.Phone}
	saveduser, err := h.usrsvc.CreateUser(usr)
	if err != nil {
		return c.Status(401).JSON(res)
	}
	res.Data = saveduser
	res.Success = true
	res.Error = ""
	return c.JSON(res)
}
