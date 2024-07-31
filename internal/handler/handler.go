package handler

import (
	"encoding/json"
	"net/http"

	"github.com/encountea/message-service/internal/models"
	"github.com/encountea/message-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// func InitRoutes(r *http.ServeMux, h *Handler) *http.ServeMux {
// 	r.HandleFunc("/messages", h.CreateMessage)
// 	r.HandleFunc("/stats", h.GetStats)
// 	return r
// }

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/messages", h.CreateMessage)
	mux.HandleFunc("/stats", h.GetStats)
	return mux
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ProcessMessage(msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.service.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}