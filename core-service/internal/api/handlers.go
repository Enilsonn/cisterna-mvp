package api

import (
	"cisterna-mvp/core-service/internal/repository"
	"cisterna-mvp/core-service/internal/utils"
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

func (h *mapHandler) HandlerGetTruckCurrentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "truckID")
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
