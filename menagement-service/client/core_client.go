package client

import (
	"bytes"
	"cisterna-mvp/menagement-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CoreClient interface {
	SyncCistern(ctx context.Context, c domain.Cistern) error
}

type coreClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewCoreClient(baseURL string) CoreClient {
	return &coreClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *coreClientImpl) SyncCistern(ctx context.Context, cistern domain.Cistern) error {
	payload := map[string]any{
		"id":              cistern.ID,
		"name":            cistern.Name,
		"capacity_liters": cistern.CapacityLiters,
		"latitude":        cistern.Latitude,
		"longitude":       cistern.Longitude,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error to marshal payload json: %v", err)
	}

	url := fmt.Sprintf("%s/api/v1/cisterns", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error to open new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error to invoke core-service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("core-service returns an unexpectable error: %v", err)
	}

	return nil
}
