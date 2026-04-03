package api

import (
	"cisterna-mvp/menagement-service/client"
	"cisterna-mvp/menagement-service/internal/domain"
	"cisterna-mvp/menagement-service/internal/repository"
	"cisterna-mvp/menagement-service/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiHandler struct {
	repo       repository.SighRepository
	coreClient client.CoreClient
}

func NewApiHandler(repo repository.SighRepository, coreClient client.CoreClient) *apiHandler {
	return &apiHandler{
		repo:       repo,
		coreClient: coreClient,
	}
}

func (h *apiHandler) CreateCistern(w http.ResponseWriter, r *http.Request) {
	var cistern domain.Cistern
	if err := json.NewDecoder(r.Body).Decode(&cistern); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	id, err := h.repo.CreateCistern(r.Context(), cistern)
	if err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to create cistern in handler: %v", err))
		return
	}

	cistern.ID = id

	if err := h.coreClient.SyncCistern(r.Context(), cistern); err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to sync cistern in core-service: %v", err))
	}
	utils.RespondeWithJSON(w, http.StatusCreated, cistern)
}

func (h *apiHandler) CreatePipeiro(w http.ResponseWriter, r *http.Request) {
	var pipeiro domain.Pipeiro
	if err := json.NewDecoder(r.Body).Decode(&pipeiro); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	id, err := h.repo.CreatePipeiro(r.Context(), pipeiro)
	if err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to create pipeiro in handler: %v", err))
		return
	}

	pipeiro.ID = id
	utils.RespondeWithJSON(w, http.StatusCreated, pipeiro)
}

func (h *apiHandler) CreateTruck(w http.ResponseWriter, r *http.Request) {
	var truck domain.Truck
	if err := json.NewDecoder(r.Body).Decode(&truck); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	id, err := h.repo.CreateTruck(r.Context(), truck)
	if err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to create truck in handler: %v", err))
		return
	}

	truck.ID = id
	utils.RespondeWithJSON(w, http.StatusCreated, truck)
}

func (h *apiHandler) CreateDelivery(w http.ResponseWriter, r *http.Request) {
	var delivery domain.Delivery
	if err := json.NewDecoder(r.Body).Decode(&delivery); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	id, err := h.repo.CreateDelivery(r.Context(), delivery)
	if err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to create delivery in handler: %v", err))
		return
	}

	delivery.ID = id
	utils.RespondeWithJSON(w, http.StatusCreated, delivery)
}

func (h *apiHandler) UpdateCistern(w http.ResponseWriter, r *http.Request) {
	var cistern domain.Cistern
	if err := json.NewDecoder(r.Body).Decode(&cistern); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	if err := h.repo.UpdateCistern(r.Context(), cistern); err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to update cistern in handler: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusCreated, cistern)
}

func (h *apiHandler) UpdateTruck(w http.ResponseWriter, r *http.Request) {
	var truck domain.Truck
	if err := json.NewDecoder(r.Body).Decode(&truck); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	if err := h.repo.UpdateTruck(r.Context(), truck); err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to update truck in handler: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusCreated, truck)
}

func (h *apiHandler) UpdatePipeiro(w http.ResponseWriter, r *http.Request) {
	var pipeiro domain.Pipeiro
	if err := json.NewDecoder(r.Body).Decode(&pipeiro); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	if err := h.repo.UpdatePipeiro(r.Context(), pipeiro); err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to update pipeiro in handler: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusCreated, pipeiro)
}

func (h *apiHandler) UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	var delivery domain.Delivery
	if err := json.NewDecoder(r.Body).Decode(&delivery); err != nil {
		utils.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("error to decode JSON: %v", err))
		return
	}

	if err := h.repo.UpdateDelivery(r.Context(), delivery); err != nil {
		utils.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("error to update delivery in handler: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusCreated, delivery)
}

func (h *apiHandler) GetPipeiroByCPF(w http.ResponseWriter, r *http.Request) {
	cpf := chi.URLParam(r, "cpf")

	if cpf == "" || len(cpf) < 11 || len(cpf) > 14 {
		utils.RespondeWithError(w, http.StatusBadRequest, "invalid cpf")
		return
	}

	pipeiro, err := h.repo.GetPipeiroByCPF(r.Context(), cpf)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get pipeiro by cpf: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusOK, pipeiro)
}

func (h *apiHandler) GetTruckByPlate(w http.ResponseWriter, r *http.Request) {
	plate := chi.URLParam(r, "plate")

	if plate == "" || len(plate) > 15 {
		utils.RespondeWithError(w, http.StatusBadRequest, "invalid plate")
		return
	}

	truck, err := h.repo.GetTruckByPlate(r.Context(), plate)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get truck by plate: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusOK, truck)
}

func (h *apiHandler) GetCisternByUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	if uuid == "" {
		utils.RespondeWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	cistern, err := h.repo.GetCisternByUUID(r.Context(), uuid)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get cistern by uuid: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusOK, cistern)
}

func (h *apiHandler) GetDeliveryByUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	if uuid == "" {
		utils.RespondeWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	delivery, err := h.repo.GetDeliveryByUUID(r.Context(), uuid)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get delivery by uuid: %v", err))
		return
	}

	utils.RespondeWithJSON(w, http.StatusOK, delivery)
}

func (h *apiHandler) GetTruckByPipeiroUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		utils.RespondeWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	trucks, err := h.repo.GetTruckByPipeiroUUID(r.Context(), uuid)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get truck by pipeiro uuid: %v", err))
		return
	}

	if trucks == nil {
		trucks = []*domain.Truck{}
	}

	utils.RespondeWithJSON(w, http.StatusOK, trucks)
}

func (h *apiHandler) GetCisterns(w http.ResponseWriter, r *http.Request) {
	cisterns, err := h.repo.GetCisterns(r.Context())
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get cisterns: %v", err))
		return
	}

	if cisterns == nil {
		cisterns = []*domain.Cistern{}
	}

	utils.RespondeWithJSON(w, http.StatusOK, cisterns)
}

func (h *apiHandler) GetDeliveryByPipeiroUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		utils.RespondeWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	deliveries, err := h.repo.GetDeliveryByPipeiroUUID(r.Context(), uuid)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get deliveries by pipeiro uuid: %v", err))
		return
	}

	if deliveries == nil {
		deliveries = []*domain.Delivery{}
	}

	utils.RespondeWithJSON(w, http.StatusOK, deliveries)
}

func (h *apiHandler) GetDeliveryByTruckUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		utils.RespondeWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	deliveries, err := h.repo.GetDeliveryByTruckUUID(r.Context(), uuid)
	if err != nil {
		utils.RespondeWithError(w, http.StatusNotFound, fmt.Sprintf("error to get deliveries by pipeiro uuid: %v", err))
		return
	}

	if deliveries == nil {
		deliveries = []*domain.Delivery{}
	}

	utils.RespondeWithJSON(w, http.StatusOK, deliveries)
}
