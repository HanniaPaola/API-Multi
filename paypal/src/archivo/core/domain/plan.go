package domain

type PlanType string

const (
    Basic    PlanType = "basic"
    BasicPlus PlanType = "basic+"
)

type Station struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Latitude  string   `json:"latitude"`
	Longitude string   `json:"longitude"`
	OwnerID   int      `json:"owner_id"`
	Plan      PlanType `json:"plan"`
}

type PlanDetails struct {
	Type    PlanType `json:"type"`
	Title   string   `json:"title"`
	Price   float64  `json:"price"`
	Active  bool     `json:"active"`
}

var PlansConfig = map[PlanType]PlanDetails{
	Basic: {
		Type:   Basic,
		Title:  "Plan Básico",
		Price:  399.00,
		Active: true,
	},
	BasicPlus: {
		Type:   BasicPlus,
		Title:  "Plan Básico Plus",
		Price:  899.00,
		Active: true,
	},
}