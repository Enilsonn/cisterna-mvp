package repository

import (
	"cisterna-mvp/menagement-service/internal/domain"
	"context"
	"fmt"
	"sync"
	"time"
)

type inMemorySighRepo struct {
	mu sync.RWMutex

	pipeiros   map[string]domain.Pipeiro
	trucks     map[string]domain.Truck
	cisterns   map[string]domain.Cistern
	deliveries map[string]domain.Delivery
}

func NewInMemorySighRepo() SighRepository {
	return &inMemorySighRepo{
		pipeiros:   make(map[string]domain.Pipeiro),
		trucks:     make(map[string]domain.Truck),
		cisterns:   make(map[string]domain.Cistern),
		deliveries: make(map[string]domain.Delivery),
	}
}

// ─── Pipeiro ────────────────────────────────────────────────────────────────

func (r *inMemorySighRepo) CreatePipeiro(_ context.Context, pipeiro domain.Pipeiro) (string, error) {
	if pipeiro.ID == "" {
		return "", fmt.Errorf("pipeiro ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pipeiros[pipeiro.ID]; exists {
		return "", fmt.Errorf("pipeiro(id: %v) already exists", pipeiro.ID)
	}

	pipeiro.CreatedAt = time.Now()
	r.pipeiros[pipeiro.ID] = pipeiro
	return pipeiro.ID, nil
}

func (r *inMemorySighRepo) UpdatePipeiro(_ context.Context, pipeiro domain.Pipeiro) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.pipeiros[pipeiro.ID]
	if !ok {
		return fmt.Errorf("no one row affected")
	}

	// preserva campos imutáveis
	pipeiro.CreatedAt = existing.CreatedAt
	r.pipeiros[pipeiro.ID] = pipeiro
	return nil
}

func (r *inMemorySighRepo) GetPipeiroByCPF(_ context.Context, cpf string) (*domain.Pipeiro, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.pipeiros {
		if p.CPF == cpf {
			cp := p
			return &cp, nil
		}
	}

	return nil, fmt.Errorf("pipeiro(cpf: %v) not found", cpf)
}

// ─── Truck ──────────────────────────────────────────────────────────────────

func (r *inMemorySighRepo) CreateTruck(_ context.Context, truck domain.Truck) (string, error) {
	if truck.ID == "" {
		return "", fmt.Errorf("truck ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.trucks[truck.ID]; exists {
		return "", fmt.Errorf("truck(id: %v) already exists", truck.ID)
	}

	truck.CreatedAt = time.Now()
	r.trucks[truck.ID] = truck
	return truck.ID, nil
}

func (r *inMemorySighRepo) UpdateTruck(_ context.Context, truck domain.Truck) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.trucks[truck.ID]
	if !ok {
		return fmt.Errorf("no one row affected")
	}

	truck.CreatedAt = existing.CreatedAt
	r.trucks[truck.ID] = truck
	return nil
}

func (r *inMemorySighRepo) GetTruckByPlate(_ context.Context, plate string) (*domain.Truck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, t := range r.trucks {
		if t.Plate == plate {
			cp := t
			return &cp, nil
		}
	}

	return nil, fmt.Errorf("truck(plate: %v) not found", plate)
}

func (r *inMemorySighRepo) GetTruckByPipeiroUUID(_ context.Context, uuid string) ([]*domain.Truck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Truck
	for _, t := range r.trucks {
		if t.PipeiroID == uuid {
			cp := t
			result = append(result, &cp)
		}
	}

	return result, nil
}

// ─── Cistern ────────────────────────────────────────────────────────────────

func (r *inMemorySighRepo) CreateCistern(_ context.Context, cistern domain.Cistern) (string, error) {
	if cistern.ID == "" {
		return "", fmt.Errorf("cistern ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cisterns[cistern.ID]; exists {
		return "", fmt.Errorf("cistern(id: %v) already exists", cistern.ID)
	}

	cistern.CreatedAt = time.Now()
	r.cisterns[cistern.ID] = cistern
	return cistern.ID, nil
}

func (r *inMemorySighRepo) UpdateCistern(_ context.Context, cistern domain.Cistern) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.cisterns[cistern.ID]
	if !ok {
		return fmt.Errorf("no one row affected")
	}

	cistern.CreatedAt = existing.CreatedAt
	r.cisterns[cistern.ID] = cistern
	return nil
}

func (r *inMemorySighRepo) GetCisternByUUID(_ context.Context, uuid string) (*domain.Cistern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.cisterns[uuid]
	if !ok {
		return nil, fmt.Errorf("cistern not found")
	}

	cp := c
	return &cp, nil
}

func (r *inMemorySighRepo) GetCisterns(_ context.Context) ([]*domain.Cistern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Cistern, 0, len(r.cisterns))
	for _, c := range r.cisterns {
		cp := c
		result = append(result, &cp)
	}

	return result, nil
}

// ─── Delivery ────────────────────────────────────────────────────────────────

func (r *inMemorySighRepo) CreateDelivery(_ context.Context, delivery domain.Delivery) (string, error) {
	if delivery.ID == "" {
		return "", fmt.Errorf("delivery ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.deliveries[delivery.ID]; exists {
		return "", fmt.Errorf("delivery(id: %v) already exists", delivery.ID)
	}

	now := time.Now()
	delivery.CreatedAt = now
	delivery.UpdatedAt = now
	if delivery.Status == "" {
		delivery.Status = domain.StatusSchelued
	}

	r.deliveries[delivery.ID] = delivery
	return delivery.ID, nil
}

func (r *inMemorySighRepo) UpdateDelivery(_ context.Context, delivery domain.Delivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.deliveries[delivery.ID]
	if !ok {
		return fmt.Errorf("no one row affected")
	}

	delivery.CreatedAt = existing.CreatedAt
	delivery.UpdatedAt = time.Now()
	r.deliveries[delivery.ID] = delivery
	return nil
}

func (r *inMemorySighRepo) GetDeliveryByUUID(_ context.Context, uuid string) (*domain.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.deliveries[uuid]
	if !ok {
		return nil, fmt.Errorf("delivery not found")
	}

	cp := d
	return &cp, nil
}

func (r *inMemorySighRepo) GetDeliveryByPipeiroUUID(_ context.Context, uuid string) ([]*domain.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// replica o JOIN: deliveries → trucks → pipeiro_id
	var result []*domain.Delivery
	for _, d := range r.deliveries {
		truck, ok := r.trucks[d.TruckID]
		if !ok || truck.PipeiroID != uuid {
			continue
		}
		cp := d
		result = append(result, &cp)
	}

	return result, nil
}

func (r *inMemorySighRepo) GetDeliveryByTruckUUID(_ context.Context, uuid string) ([]*domain.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Delivery
	for _, d := range r.deliveries {
		if d.TruckID == uuid {
			cp := d
			result = append(result, &cp)
		}
	}

	return result, nil
}

func (r *inMemorySighRepo) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.pipeiros = make(map[string]domain.Pipeiro)
	r.trucks = make(map[string]domain.Truck)
	r.cisterns = make(map[string]domain.Cistern)
	r.deliveries = make(map[string]domain.Delivery)
}

var _ SighRepository = (*inMemorySighRepo)(nil)
