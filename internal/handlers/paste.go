package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"pastebin-backend/internal/models"
	"pastebin-backend/internal/repository"
)

type PasteHandler struct {
	repo *repository.PasteRepository
}

func NewPasteHandler(repo *repository.PasteRepository) *PasteHandler {
	return &PasteHandler{repo: repo}
}

// createPaste handles POST request to create new paste
func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var paste models.Paste
	if err := json.NewDecoder(r.Body).Decode(&paste); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(ctx, &paste); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(paste)
}

// getPaste handles GET request to retrieve paste by id
func (h *PasteHandler) GetPaste(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	paste, err := h.repo.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if paste == nil {
		http.Error(w, "paste not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paste)
}

// updatePaste handles PUT request to update existing paste
func (h *PasteHandler) UpdatePaste(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	var paste models.Paste
	if err := json.NewDecoder(r.Body).Decode(&paste); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paste.ID = id

	if err := h.repo.Update(ctx, &paste); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paste)
}

// deletePaste handles DELETE request to remove paste
func (h *PasteHandler) DeletePaste(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getAllPastes handles GET request to retrieve all pastes
func (h *PasteHandler) GetAllPastes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pastes, err := h.repo.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pastes)
}
