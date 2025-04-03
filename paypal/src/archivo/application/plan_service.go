package application

import (
	"api/paypal/src/archivo/core/domain"
	"api/paypal/src/archivo/core/ports"
)

type PlanService struct {
	repo ports.PlanRepository
}

func NewPlanService(repo ports.PlanRepository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) GetAvailablePlans() ([]domain.PlanDetails, error) {
	var availablePlans []domain.PlanDetails
	for _, plan := range domain.PlansConfig {
		if plan.Active {
			availablePlans = append(availablePlans, plan)
		}
	}
	return availablePlans, nil
}

func (s *PlanService) GetStationPlan(stationID int) (*domain.Station, *domain.PlanDetails, error) {
	station, err := s.repo.GetStationByID(stationID)
	if err != nil {
		return nil, nil, err
	}

	planDetails, err := s.repo.GetPlanDetails(station.Plan)
	if err != nil {
		return nil, nil, err
	}

	return station, planDetails, nil
}

func (s *PlanService) UpgradePlan(stationID int, newPlan domain.PlanType) (*domain.Station, error) {
	return s.repo.UpdateStationPlan(stationID, newPlan)
}
