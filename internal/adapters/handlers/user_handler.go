package handler

import (
	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usrsvc   ports.UserService
	authsvc  ports.AuthService
	emailsvc ports.EmailService
	orgsvc   ports.OrgService
}

func NewUserHandler(svc ports.UserService, authsvc ports.AuthService, emailsvc ports.EmailService, orgsvc ports.OrgService) *UserHandler {
	return &UserHandler{usrsvc: svc, authsvc: authsvc, emailsvc: emailsvc, orgsvc: orgsvc}
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
	res := domain.Envelope[domain.User]{Success: false, Data: domain.User{}, Error: ""}
	if err := c.BodyParser(&req); err != nil {
		res.Error = "Bad request"
		return c.Status(400).JSON(res)
	}
	pword, err := h.authsvc.HashPassword(req.Password)
	usr := domain.User{Name: req.Name, Email: req.Email, Password: pword, OrgID: req.OrgID, Role: req.Role, Status: "pending", Phone: req.Phone}
	exists, err := h.usrsvc.IsExistingUser(usr.Email)
	if exists == true {
		err = h.usrsvc.RemoveByEmail(usr.Email)
		err = h.authsvc.RemoveByEmail(usr.Email)
	}
	orgexists, err := h.orgsvc.OrgExistsByID(usr.OrgID)
	if orgexists == false {
		res.Error = "No Organization found"
		return c.Status(400).JSON(res)
	}

	savedusr, err := h.usrsvc.CreateUser(usr)
	ent := domain.Auth{Email: savedusr.Email, UserID: savedusr.ID}
	resetent, _ := h.authsvc.Reset(ent, "48h")
	emailBody := h.emailsvc.GetWelcomeEmailBody(resetent.Token)
	err = h.emailsvc.SendEmail(req.Email, "[Action Required] Welcome; Martian 🧑‍🚀.", emailBody)
	if err != nil {
		res.Success = false
		res.Error = "Error sending email"
		return c.Status(401).JSON(res)
	}
	res.Success = true
	res.Data = savedusr
	res.Error = ""
	return c.Status(200).JSON(res)
}
func (h *UserHandler) Get(c *fiber.Ctx) error {
	res := domain.Envelope[[]domain.User]{Success: false, Data: []domain.User{}, Error: ""}
	queryParams := c.Queries()
	org_id := queryParams["org_id"]
	allusers, err := h.usrsvc.FindAllByOrgID(org_id)
	if err != nil {
		res.Error = "Could not find users for this org"
		return c.Status(400).JSON(res)
	}
	res.Success = true
	res.Error = ""
	res.Data = allusers
	return c.Status(200).JSON(res)
}
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Name     string          `json:"name"`
		Email    string          `json:"email"`
		Password string          `json:"password"`
		OrgID    string          `json:"org_id"`
		Role     domain.Roles    `json:"role"`
		Status   domain.Statuses `json:"status"`
		Phone    string          `json:"phone"`
	}
	res := domain.Envelope[domain.User]{Success: false, Data: domain.User{}, Error: ""}
	if err := c.BodyParser(&req); err != nil {
		res.Error = "Bad request"
		return c.Status(400).JSON(res)
	}
	pword, err := h.usrsvc.FindPasswordByEmail(req.Email)
	correctpword := h.authsvc.CheckPasswordHash(req.Password, pword)
	if err != nil || correctpword == false {
		res.Error = "Incorrect Password"
		return c.Status(400).JSON(res)
	}
	profile, err := h.usrsvc.FindByEmail(req.Email)
	if err != nil {
		res.Success = false
		res.Error = "Error getting User Profile"
		return c.Status(400).JSON(res)
	}
	res.Success = true
	res.Data = profile
	res.Error = ""
	return c.Status(200).JSON(res)
}
func (h *UserHandler) Remove(c *fiber.Ctx) error {
	var req struct {
		Email string       `json:"email"`
		OrgID string       `json:"org_id"`
		Role  domain.Roles `json:"role"`
	}
	res := domain.Envelope[domain.User]{Success: false, Data: domain.User{}, Error: ""}
	if err := c.BodyParser(&req); err != nil {
		res.Error = "Bad request"
		return c.Status(400).JSON(res)
	}
	if req.Role != "admin" {
		res.Error = "Role type member are not allowed to perform this action. Contact admin"
		return c.Status(400).JSON(res)
	}
	profile, err := h.usrsvc.FindByEmail(req.Email)
	if err != nil {
		res.Error = "No user with that email"
		return c.Status(400).JSON(res)
	}
	err = h.usrsvc.RemoveByEmail(profile.Email)
	if err != nil {
		res.Success = false
		res.Error = "Error sending email"
		return c.Status(401).JSON(res)
	}
	res.Success = true
	res.Error = ""
	return c.Status(200).JSON(res)
}
