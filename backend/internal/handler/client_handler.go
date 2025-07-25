package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"controle-de-estoque/backend/internal/domain"
	"controle-de-estoque/backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ClientHandler gerencia as requisições HTTP para clientes.
type ClientHandler struct {
	service *service.ClientService
}

// NewClientHandler cria uma nova instância de ClientHandler.
func NewClientHandler(s *service.ClientService) *ClientHandler {
	return &ClientHandler{service: s}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client domain.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}
	if err := h.service.Create(r.Context(), &client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(client); err != nil {
		log.Printf("Erro ao codificar JSON do cliente: %v", err)
	}
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(clients); err != nil {
		log.Printf("Erro ao codificar JSON da lista de clientes: %v", err)
	}
}

func (h *ClientHandler) GetClientByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "clientID")
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
		return
	}
	client, err := h.service.GetByID(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(client); err != nil {
		log.Printf("Erro ao codificar JSON do cliente: %v", err)
	}
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "clientID")
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
		return
	}
	var client domain.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}
	client.ID = clientID
	if err := h.service.Update(r.Context(), &client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(client); err != nil {
		log.Printf("Erro ao codificar JSON do cliente atualizado: %v", err)
	}
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "clientID")
	clientID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
		return
	}
	if err := h.service.Delete(r.Context(), clientID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
