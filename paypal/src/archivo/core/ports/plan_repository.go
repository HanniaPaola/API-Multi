package ports

import "api/paypal/src/archivo/core/domain"

type PlanRepository interface {
	GetStationByID(id int) (*domain.Station, error)
	UpdateStationPlan(stationID int, newPlan domain.PlanType) (*domain.Station, error)
	GetPlanDetails(planType domain.PlanType) (*domain.PlanDetails, error)
}
