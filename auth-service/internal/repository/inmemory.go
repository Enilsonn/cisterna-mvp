package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"cisterna-mvp/auth-service/internal/domain"
)

type InMemoryRepository struct {
	mu      sync.RWMutex
	users   map[string]domain.User
	refresh map[string]domain.RefreshToken
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users:   make(map[string]domain.User),
		refresh: make(map[string]domain.RefreshToken),
	}
}

func (r *InMemoryRepository) CreateUser(ctx context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, errors.New("user not found")
}

func (r *InMemoryRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryRepository) SaveRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.refresh[token.Token] = token
	return nil
}

func (r *InMemoryRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.refresh[token]
	if !ok {
		return fmt.Errorf("refresh token not found")
	}
	entry.Revoked = true
	r.refresh[token] = entry
	return nil
}

func (r *InMemoryRepository) GetRefreshToken(ctx context.Context, token string) (domain.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.refresh[token]
	if !ok {
		return domain.RefreshToken{}, fmt.Errorf("refresh token not found")
	}
	return entry, nil
}
