package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type ExistsResponse struct {
	Exists bool `json:"exists"`
}

type ResolverService struct {
	cfg    *Config
	client *http.Client
}

func NewResolverService(cfg *Config) *ResolverService {
	return &ResolverService{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.ExistsTimeoutMs) * time.Millisecond,
		},
	}
}

func (s *ResolverService) Resolve(ctx context.Context, imageURL string) (string, error) {
	imageURL = strings.TrimSpace(imageURL)

	idx := strings.Index(imageURL, ".ru/")
	if idx == -1 {
		return "", fmt.Errorf("invalid url")
	}

	path := strings.TrimSpace(imageURL[idx+3:])

	existsURL := fmt.Sprintf("%s/exists%s", s.cfg.OurDomain, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, existsURL, nil)
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return s.cfg.GoldAppleDomain + path, nil
	}

	resp, err := s.client.Do(req)
	if err != nil {
		slog.Error("request failed", "error", err)
		return s.cfg.GoldAppleDomain + path, nil
	}
	defer resp.Body.Close()

	var existsResp ExistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&existsResp); err != nil {
		slog.Error("failed to decode response", "error", err)
		return s.cfg.GoldAppleDomain + path, nil
	}

	if existsResp.Exists {
		return s.cfg.OurDomain + path, nil
	}

	return s.cfg.GoldAppleDomain + path, nil
}
