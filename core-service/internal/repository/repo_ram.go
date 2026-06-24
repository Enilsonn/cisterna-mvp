package repository

import (
	"cisterna-mvp/core-service/internal/domain"
	"context"
	"fmt"
	"sync"
)

type inMemoryRepo struct {
	mu sync.RWMutex

	positionHistory map[string][]domain.TruckPosition
	currentStatus   map[string]domain.TruckStatus
	cisterns        map[int]domain.Cisterna

	nextCisternID int
}

func NewInMemoryRepo() PositionRepository {
	return &inMemoryRepo{
		positionHistory: make(map[string][]domain.TruckPosition),
		currentStatus:   make(map[string]domain.TruckStatus),
		cisterns:        make(map[int]domain.Cisterna),
		nextCisternID:   1,
	}
}

func (r *inMemoryRepo) SavePosition(_ context.Context, pos domain.TruckPosition) error {
	if pos.TruckID == "" {
		return fmt.Errorf("truck_id cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// INSERT
	r.positionHistory[pos.TruckID] = append(r.positionHistory[pos.TruckID], pos)

	// UPSERT
	r.currentStatus[pos.TruckID] = domain.TruckStatus{
		TruckID:   pos.TruckID,
		Longitude: fmt.Sprintf("%f", pos.Longitude),
		Latitude:  pos.Latitude,
		LastSeen:  pos.Timestamp,
	}

	return nil
}

func (r *inMemoryRepo) GetTruckCurrrentLocation(_ context.Context, id string) (*domain.TruckStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	status, ok := r.currentStatus[id]
	if !ok {
		return nil, fmt.Errorf("truck not found: %s", id)
	}

	return &status, nil
}

func (r *inMemoryRepo) CreateCistern(_ context.Context, cis domain.Cisterna) (int64, error) {
	if cis.Nome == "" {
		return 0, fmt.Errorf("cistern name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	cis.ID = r.nextCisternID
	r.cisterns[r.nextCisternID] = cis
	r.nextCisternID++

	return int64(cis.ID), nil
}

func (r *inMemoryRepo) GetCisterns(_ context.Context) ([]domain.Cisterna, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Cisterna, 0, len(r.cisterns))
	for _, c := range r.cisterns {
		result = append(result, c)
	}

	return result, nil
}

// funcoes que ainda não tem na inferface mas podem ser acessadas pela iniciação nominal da struct
func (r *inMemoryRepo) GetPositionHistory(truckID string) []domain.TruckPosition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	history := r.positionHistory[truckID]
	// retorna cópia para evitar race condition fora do lock
	cp := make([]domain.TruckPosition, len(history))
	copy(cp, history)
	return cp
}

// limpar a cash -- perder os ponteiros antigos
func (r *inMemoryRepo) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.positionHistory = make(map[string][]domain.TruckPosition)
	r.currentStatus = make(map[string]domain.TruckStatus)
	r.cisterns = make(map[int]domain.Cisterna)
	r.nextCisternID = 1
}

// professor raoni, essa linha abaixo serve para forcar o compilador verificar em tempo de compilacao
// que a struct inMemoryRepo tem os métodos da interface PositionRepository, o que torna bem mais rapido a criação da "classe"
var _ PositionRepository = (*inMemoryRepo)(nil)
