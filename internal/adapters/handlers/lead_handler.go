package handler

import (
	"fmt"
	"time"

	"github.com/strconvitoa/martian-service/internal/core/domain"
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

type LeadHandler struct {
	leadsvc ports.LeadService
}

func NewLeadHandler(leadsvc ports.LeadService) *LeadHandler {
	return &LeadHandler{leadsvc: leadsvc}
}

func (h *LeadHandler) Create(c *fiber.Ctx) error {
	req := domain.Lead{}
	res := domain.Envelope[domain.Lead]{Success: false, Data: domain.Lead{}, Error: ""}
	if err := c.BodyParser(&req); err != nil {
		res.Error = "Bad Request"
		return c.Status(400).JSON(res)
	}
	Lead, err := h.leadsvc.CreateLead(req)
	if err != nil {
		res.Success = false
		res.Data = Lead
		res.Error = "Error creating Lead"
		return c.Status(401).JSON(res)
	}
	res.Success = true
	res.Data = Lead
	res.Error = ""
	return c.Status(200).JSON(res)
}
func (h *LeadHandler) Update(c *fiber.Ctx) error {
	req := domain.Lead{}
	res := domain.Envelope[domain.Lead]{Success: false, Data: domain.Lead{}, Error: ""}
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing ID"})
	}
	if err := c.BodyParser(&req); err != nil {
		res.Error = "Bad Request"
		return c.Status(400).JSON(res)
	}
	lead, err := h.leadsvc.ChageLeadStatus(id, req.Status)
	if err != nil {
		res.Success = false
		res.Error = "Error creating Lead"
		return c.Status(401).JSON(res)
	}
	res.Success = true
	res.Data = lead
	res.Error = ""
	return c.Status(200).JSON(res)
}

func (h *LeadHandler) Get(c *fiber.Ctx) error {

	res := domain.Envelope[domain.LeadsResp]{Success: false, Data: domain.LeadsResp{}, Error: ""}
	queryParams := c.Queries()
	org_id := queryParams["org_id"]
	dataPayload := domain.LeadsResp{Leads: []domain.Lead{}, Discarded: []domain.Lead{}, Pending: []domain.Lead{}, NewLeads: 0}
	aleads, err := h.leadsvc.FindLeadByStatus(org_id, "accept")
	dleads, err := h.leadsvc.FindLeadByStatus(org_id, "discard")
	pleads, err := h.leadsvc.FindLeadByStatus(org_id, "pending")

	// Format Create date
	pleads, err = h.FormatLeadsToFriendlyDate(pleads)
	aleads, err = h.FormatLeadsToFriendlyDate(aleads)
	dleads, err = h.FormatLeadsToFriendlyDate(dleads)
	if err != nil {
		res.Error = "issue finding leads"
		return c.Status(400).JSON(res)
	}
	dataPayload.Leads = aleads
	dataPayload.Discarded = dleads
	dataPayload.Pending = pleads

	res.Success = true
	res.Data = dataPayload
	res.Error = ""
	return c.Status(200).JSON(res)
}

func (h *LeadHandler) FormatLeadsToFriendlyDate(leads []domain.Lead) ([]domain.Lead, error) {
	// The layout matching your database string format
	const dbLayout = "2006-01-02 15:04:05.999999-07"

	// The UI friendly layout requested: "January 2, 2006"
	const friendlyLayout = "January 2, 2006"

	// Use an index loop to mutate the elements directly in the slice
	for i := 0; i < len(leads); i++ {
		// Skip empty strings if any exist to prevent parsing errors
		if leads[i].Created == "" {
			continue
		}

		// Parse the database timestamp string
		parsedTime, err := time.Parse(dbLayout, leads[i].Created)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date '%s' for lead ID %s: %w", leads[i].Created, leads[i].ID, err)
		}

		// Overwrite the Created property with the friendly format
		leads[i].Created = parsedTime.Format(friendlyLayout)
	}

	return leads, nil
}
