package domain

import "time"

type TruckPosition struct {
	TruckID   string
	Longitude float64
	Latitude  float64
	Timestamp time.Time
}
