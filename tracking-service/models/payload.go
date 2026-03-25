package models

// Representa o "envelope" que o app do caminhão irá enviar
type GPSPayload struct {
	TruckID   string  `json:"truck_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}
