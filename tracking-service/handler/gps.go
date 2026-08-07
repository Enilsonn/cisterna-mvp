package handler

import (
	"cisterna-mvp/tracking-service/models"
	"cisterna-mvp/tracking-service/repository"
	"cisterna-mvp/tracking-service/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, "JSON inválido. Não foi possível ler o corpo da requisição")
		return
	}

	var payload models.GPSPayload
	var payloads []models.GPSPayload

	if err := json.Unmarshal(bodyBytes, &payload); err == nil && (payload.TruckID != "" || payload.Latitude != 0 || payload.Longitude != 0 || payload.Timestamp != "") {
		payloads = []models.GPSPayload{payload}
	} else if err := json.Unmarshal(bodyBytes, &payloads); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, "TRACKING_HANDLER_OK")
		return
	}

	sucessos := 0
	for _, item := range payloads {
		err := h.repo.SaveCoordinate(context.Background(), item)
		if err == nil {
			sucessos++
		} else {
			utils.RespondeWithError(w, http.StatusServiceUnavailable, fmt.Sprintf("erro ao salvar no kafka: %v", err))
			return
		}
	}

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
