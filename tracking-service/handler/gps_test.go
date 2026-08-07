package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cisterna-mvp/tracking-service/models"
)

type stubTrackingRepo struct {
	calls int
}

func (s *stubTrackingRepo) SaveCoordinate(_ context.Context, _ models.GPSPayload) error {
	s.calls++
	return nil
}

func TestReciveGPSAcceptsSinglePayload(t *testing.T) {
	repo := &stubTrackingRepo{}
	h := NewGPSHandler(repo)

	body := `{"truck_id":"truck-1","latitude":-22.97,"longitude":-43.21,"timestamp":"2026-08-07T00:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gps", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ReciveGPS(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}

	if repo.calls != 1 {
		t.Fatalf("expected 1 saved coordinate, got %d", repo.calls)
	}
}
