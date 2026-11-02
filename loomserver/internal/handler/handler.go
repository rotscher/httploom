// Package handler implements the generated Loom API server interface.
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rotscher/httploom/loomserver/internal/api"
	"github.com/rotscher/httploom/loomserver/internal/weave"
)

const maxBodyBytes = 1 << 20

// Handler serves the Loom API backed by an in-memory store.
type Handler struct {
	store   *weave.Store
	version string
}

// New returns a Handler that reports the given version.
func New(store *weave.Store, version string) *Handler {
	return &Handler{store: store, version: version}
}

// GetHealth implements GET /health.
func (h *Handler) GetHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, api.Health{Status: "ok", Version: h.version})
}

// ListWeaves implements GET /weaves.
func (h *Handler) ListWeaves(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, api.WeaveList{Items: h.store.List()})
}

// CreateWeave implements POST /weaves.
func (h *Handler) CreateWeave(w http.ResponseWriter, r *http.Request) {
	var req api.CreateWeaveJSONRequestBody

	body := http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		writeError(w, http.StatusBadRequest, "text must not be empty")
		return
	}

	width := weave.DefaultWidth
	if req.Width != nil {
		width = *req.Width
	}
	if width < 1 || width > 200 {
		writeError(w, http.StatusBadRequest, "width must be between 1 and 200")
		return
	}

	writeJSON(w, http.StatusCreated, h.store.Create(req.Text, width))
}

// GetWeave implements GET /weaves/{id}.
func (h *Handler) GetWeave(w http.ResponseWriter, _ *http.Request, id string) {
	item, ok := h.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "weave not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, api.Error{Message: message})
}
