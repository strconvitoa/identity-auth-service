package services

import (
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/google/uuid"
	"github.com/strconvitoa/martian-service/internal/core/domain"
)

type leadService struct {
	repo ports.LeadRepository
}

func NewLeadService(repo ports.LeadRepository) ports.LeadService {
	return &leadService{repo: repo}
}

func (s *leadService) CreateLead(Leads domain.Lead) (domain.Lead, error) {
	id := uuid.New().String()
	Leads.ID = id
	savedLeads, err := s.repo.Save(Leads)
	if err != nil {
		return domain.Lead{}, err
	}
	return savedLeads, nil
}

func (s *leadService) FindLeadByStatus(org_id string, status string) ([]domain.Lead, error) {

	sleads, err := s.repo.SelectByStatus(org_id, status)
	if err != nil {
		return []domain.Lead{}, err
	}
	return sleads, nil
}

func (s *leadService) ChageLeadStatus(org_id string, status string) (domain.Lead, error) {
	sleads, err := s.repo.UpdateLeadStatus(org_id, status)
	if err != nil {
		return domain.Lead{}, err
	}
	return sleads, nil
}
