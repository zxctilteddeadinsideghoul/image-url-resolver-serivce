package main

import (
	"encoding/json"
	"net/http"
)

type ResolveRequest struct {
	ImageURL string `json:"ImageUrl"`
}

type ResolveResponse struct {
	URL string `json:"url"`
}

type Handler struct {
	service *ResolverService
}

func NewHandler(service *ResolverService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	var req ResolveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result, err := h.service.Resolve(r.Context(), req.ImageURL)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := ResolveResponse{URL: result}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
