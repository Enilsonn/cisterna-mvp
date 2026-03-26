package domain

import "time"

type TruckPosition struct {
	TruckID   string
	Longitude float64
	Latitude  float64
	Timestamp time.Time
}

type TruckStatus struct {
	TruckID   string    `json:"truck_id"`
	Longitude string    `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	LastSeen  time.Time `json:"last_seen"`
}

type Cisterna struct {
	ID             int     `json:"id"`
	Nome           string  `json:"nome"`
	CapacityLiters int     `json:"capacity_liters"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
}
