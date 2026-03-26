package api

import (
	"cisterna-mvp/core-service/internal/domain"
	"cisterna-mvp/core-service/internal/repository"
	"cisterna-mvp/core-service/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type mapHandler struct {
	repo repository.PositionRepository
}

func NewMapHandler(repo repository.PositionRepository) *mapHandler {
	return &mapHandler{
		repo: repo,
	}
}

func (h *mapHandler) HandlerGetCisters(w http.ResponseWriter, r *http.Request) {
	cisterns, err := h.repo.GetCisterns(r.Context())
	if err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to search cisterns: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusOK, cisterns)
}

func (h *mapHandler) HandlerCreateCistern(w http.ResponseWriter, r *http.Request) {
	var input domain.Cisterna

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	id, err := h.repo.CreateCistern(r.Context(), input)
	if err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to create cistern: %v", err))
		return
	}

	input.ID = int(id)
	utils.RespondeWithJSON(w, http.StatusCreated, input)
}

func (h *mapHandler) HandlerGetTruckCurrentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "truck_id")
	if id == "" {
		utils.RespondeWithError(w, http.StatusBadRequest, "no one id informed")
		return
	}

	truck, err := h.repo.GetTruckCurrrentLocation(r.Context(), id)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("no one truck found: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusOK, truck)
}
