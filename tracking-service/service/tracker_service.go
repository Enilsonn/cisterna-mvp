package service

import (
	"cisterna-mvp/tracking-service/models"
	"cisterna-mvp/tracking-service/repository"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type trackerService struct {
	repo repository.TrackingRepository
}

func NewTrackRepository(repo repository.TrackingRepository) *trackerService {
	return &trackerService{repo: repo}
}

// ProcessAndSendLocation pega a posição da rede e joga na mensageria/banco
func (s *trackerService) ProcessAndSendLocation() error {
	res, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var loc models.IPLocation
	if err = json.NewDecoder(res.Body).Decode(&loc); err != nil {
		return err
	}

	// montando o payload
	p := models.GPSPayload{
		TruckID:   "", // precisa fazer requisição ao microsserviço de autenticação
		Latitude:  loc.Lat,
		Longitude: loc.Lon,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err = s.repo.SaveCoordinate(context.Background(), p); err != nil {
		return err
	}
	return nil
}
