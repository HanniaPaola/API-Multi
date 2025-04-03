package mysql

import (
	"api/paypal/src/archivo/core/domain"
	"database/sql"
	"fmt"
)

type MySQLPlanRepository struct {
	db *sql.DB
}

func NewMySQLPlanRepository(db *sql.DB) *MySQLPlanRepository {
	return &MySQLPlanRepository{db: db}
}

func (r *MySQLPlanRepository) GetStationByID(id int) (*domain.Station, error) {
	query := `SELECT id, name, latitude, longitude, owner_id, plan FROM stations WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var station domain.Station
	err := row.Scan(
		&station.ID,
		&station.Name,
		&station.Latitude,
		&station.Longitude,
		&station.OwnerID,
		&station.Plan,
	)

	if err != nil {
		return nil, err
	}

	return &station, nil
}

func (r *MySQLPlanRepository) UpdateStationPlan(stationID int, newPlan domain.PlanType) (*domain.Station, error) {
	query := `UPDATE stations SET plan = ? WHERE id = ?`
	_, err := r.db.Exec(query, newPlan, stationID)
	if err != nil {
		return nil, err
	}

	return r.GetStationByID(stationID)
}

func (r *MySQLPlanRepository) GetPlanDetails(planType domain.PlanType) (*domain.PlanDetails, error) {
	details, exists := domain.PlansConfig[planType]
	if !exists {
		return nil, fmt.Errorf("plan type not found")
	}
	return &details, nil
}
