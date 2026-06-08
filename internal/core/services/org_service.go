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

func (s *OrgService) CreateOrg(org domain.Org) (domain.Org, error) {
	id := uuid.New().String()
	org.ID = id

	savedOrg, err := s.repo.Save(org)
	if err != nil {
		return domain.Org{}, err
	}

	return savedOrg, nil
}

func (s *OrgService) OrgExists(org domain.Org) (bool, error) {
	isOrgExists, err := s.repo.Exists(org)
	if err != nil {
		return false, err
	}
	return isOrgExists, nil
}
func (s *OrgService) OrgExistsByID(id string) (bool, error) {
	isOrgExists, err := s.repo.ExistsByID(id)
	if err != nil {
		return false, err
	}
	return isOrgExists, nil
}
