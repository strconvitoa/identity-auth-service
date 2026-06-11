package handler

import (
	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authsvc  ports.AuthService
	usrsvc   ports.UserService
	emailsvc ports.EmailService
}

func NewAuthHandler(authsvc ports.AuthService, usrsvc ports.UserService, emailsvc ports.EmailService) *AuthHandler {
	return &AuthHandler{authsvc: authsvc, usrsvc: usrsvc, emailsvc: emailsvc}
}

func (h *AuthHandler) Reset(c *fiber.Ctx) error {
	var req struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		EntName   string `json:"entname"`
		TradeName string `json:"trade_name"`
		Password  string `json:"password"`
	}
	res := domain.Envelope[domain.Auth]{Success: false, Data: domain.Auth{}, Error: ""}
	if err := c.BodyParser(&req); err != nil {
		res.Success = false
		res.Error = "Bad Request"
		return c.Status(400).JSON(res)
	}
	usr, err := h.usrsvc.FindByEmail(req.Email)
	if err != nil {
		res.Success = false
		res.Error = "No User with that email"
		return c.Status(400).JSON(res)
	}
	ent := domain.Auth{Email: req.Email, UserID: usr.ID}

	savedAuth, err := h.authsvc.Reset(ent, "60m")
	if err != nil {
		return c.Status(401).JSON(res)
	}
	emailBody := h.emailsvc.GetVerificationEmailBody(savedAuth.Token)
	err = h.emailsvc.SendEmail(req.Email, "🚀[Action Required] Reset Password", emailBody)
	if err != nil {
		return c.Status(401).JSON(res)
	}
	res.Success = true
	res.Data = savedAuth
	res.Error = ""
	return c.Status(200).JSON(res)
}
func (h *AuthHandler) Change(c *fiber.Ctx) error {
	var req struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}
	res := domain.Envelope[domain.User]{Success: false, Data: domain.User{}, Error: "Bad Request"}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(res)
	}
	pword, err := h.authsvc.HashPassword(req.Password)
	if err != nil {
		res.Error = "Bad Request, Hash Error"
		return c.Status(400).JSON(res)
	}
	authReq := domain.Auth{Email: req.Email, Token: req.Token, Password: pword}
	isValid, err := h.authsvc.Validate(authReq)
	if isValid == false {
		res.Error = "Not a valid user"
		return c.Status(400).JSON(res)
	}
	usr, _ := h.usrsvc.FindByEmail(req.Email)
	usr.Password = pword
	savedusr, err := h.usrsvc.ChangePassword(usr)
	err = h.authsvc.RemoveByEmail(authReq.Email)
	if err != nil {
		res.Error = "Failed removing email"
		return c.Status(401).JSON(res)
	}
	res.Success = true
	res.Data = savedusr
	res.Error = ""
	return c.Status(200).JSON(res)
}
