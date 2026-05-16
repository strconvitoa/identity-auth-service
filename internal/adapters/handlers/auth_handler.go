package handler

import (
	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authsvc ports.AuthService
	usrsvc  ports.UserService
}

func NewAuthHandler(authsvc ports.AuthService, usrsvc ports.UserService) *AuthHandler {
	return &AuthHandler{authsvc: authsvc, usrsvc: usrsvc}
}

func (h *AuthHandler) Reset(c *fiber.Ctx) error {
	var req struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		EntName   string `json:"entname"`
		TradeName string `json:"trade_name"`
		Password  string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Bad Request")
	}
	usr, _ := h.usrsvc.GetByEmail(req.Email)
	ent := domain.Auth{Email: req.Email, UserID: usr.ID}

	savedAuth, err := h.authsvc.Reset(ent, "60m")

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"auth": savedAuth})
}
func (h *AuthHandler) Change(c *fiber.Ctx) error {
	var req struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Bad Request")
	}
	authReq := domain.Auth{Email: req.Email, ID: req.ID, Password: req.Password}
	isValid, err := h.authsvc.Validate(authReq)
	if isValid == false {
		// return empty user for now
		return c.JSON(fiber.Map{"user": domain.User{}.ToResponse()})
	}
	usr, _ := h.usrsvc.GetByEmail(req.Email)
	usr.Password = req.Password
	savedusr, err := h.usrsvc.ChangePassword(usr)
	_, err = h.authsvc.Remove(authReq.ID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error()})
	}
	resp := savedusr.ToResponse()
	return c.JSON(fiber.Map{"user": resp})
}
