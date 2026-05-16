package services

import (
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/google/uuid"
	"github.com/strconvitoa/martian-service/internal/core/domain"
)

type OrgService struct {
	repo ports.OrgRepository
}

func NewOrgService(repo ports.OrgRepository) ports.OrgService {
	return &OrgService{repo: repo}
}

func (s *OrgService) CreateOrg(Org domain.Org) (domain.Org, error) {
	id := uuid.New().String()
	Org.ID = id

	savedOrg, err := s.repo.Save(Org)
	if err != nil {
		return domain.Org{}, err
	}

	return savedOrg, nil
}
