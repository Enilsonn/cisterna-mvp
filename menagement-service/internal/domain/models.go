package domain

import "time"

type Pipeiro struct {
	ID        string    `json:"id"`
	Name      string    `json:"nome"`
	CPF       string    `json:"cpf"`
	CNH       string    `json:"cnh"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Truck struct {
	ID             string    `json:"id"`
	Plate          string    `json:"plate"`
	CapacityLiters int       `json:"capacity_liters"`
	PipeiroID      string    `json:"pipeiro_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type Cistern struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	ResponsabibleName string    `josn:"responsable_name"`
	City              string    `json:"city"`
	CapacityLiters    int       `json:"capacity_liters"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	CreatedAt         time.Time `json:"created_at"`
}

type DeliveryStatus string

const (
	StatusSchelued  DeliveryStatus = "AGENDADO"
	StatusInRoute   DeliveryStatus = "EM_ROTA"
	StatusDelivered DeliveryStatus = "CONCLUIDO"
	StatusCanceled  DeliveryStatus = "CANCELADO"
)

type Delivery struct {
	ID            string         `json:"id"`
	CisternID     string         `json:"cistern_id"`
	TruckID       string         `json:"truck_id"`
	Status        DeliveryStatus `json:"status"`
	ScheduledDate time.Time      `json:"scheduled_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
