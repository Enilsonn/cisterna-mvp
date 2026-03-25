package handler

import (
	"cisterna-mvp/tracking-service/models"
	"cisterna-mvp/tracking-service/repository"
	"cisterna-mvp/tracking-service/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GPSHandler struct {
	repo repository.TrackingRepository
}

func NewGPSHandler(repo repository.TrackingRepository) *GPSHandler {
	return &GPSHandler{repo: repo}
}

// ProcessAndSendLocation pega a posição da rede e joga na mensageria/banco
func (h *GPSHandler) ReciveGPS(w http.ResponseWriter, r *http.Request) {
	var payloads []models.GPSPayload
	if err := json.NewDecoder(r.Body).Decode(&payloads); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, "JSON inválido. Era esperado uma lista de coordenadas")
		return
	}

	sucessos := 0
	for _, payload := range payloads {
		err := h.repo.SaveCoordinate(context.Background(), payload)
		if err == nil {
			sucessos++
		} else {
			// futuramente mandar para a DLQ: serviços de logs
			utils.RespondeWithError(w, http.StatusServiceUnavailable, fmt.Sprintf("erro ao salvar no kafka: %v", err))
		}
	}

	// mandar 500 ou 503 para que uma nova tentativa de refeita
	if len(payloads) > 0 && sucessos == 0 {
		utils.RespondeWithError(w, http.StatusServiceUnavailable, "serviço indisponível,tente novamente")
		return
	}

	responseSucess := map[string]any{
		"status": "sucess",
		"saves":  sucessos,
	}

	utils.RespondeWithJSON(w, http.StatusAccepted, responseSucess)
}
