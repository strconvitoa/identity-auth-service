package services

import (
	"github.com/strconvitoa/martian-service/internal/core/ports"

	"github.com/google/uuid"
	"github.com/strconvitoa/martian-service/internal/core/domain"
)

type intakeService struct {
	repo ports.IntakeRepository
}

func NewIntakeService(repo ports.IntakeRepository) ports.IntakeService {
	return &intakeService{repo: repo}
}

func (s *intakeService) CreateIntake(Intake domain.Intake) (domain.Intake, error) {
	id := uuid.New().String()
	Intake.ID = id
	savedIntake, err := s.repo.Save(Intake)
	if err != nil {
		return domain.Intake{}, err
	}
	return savedIntake, nil
}
