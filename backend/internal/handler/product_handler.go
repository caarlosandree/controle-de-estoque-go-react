package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"controle-de-estoque/backend/internal/domain"
	"controle-de-estoque/backend/internal/repository"
	"controle-de-estoque/backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product domain.Produto
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}
	err := h.service.CreateProduct(r.Context(), &product)
	if err != nil {
		http.Error(w, "Erro ao criar o produto", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Erro ao encodar a resposta JSON: %v", err)
	}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	response, err := h.service.ListProducts(r.Context(), search, page, limit)
	if err != nil {
		http.Error(w, "Erro ao listar os produtos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao encodar a resposta JSON: %v", err)
	}
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(chi.URLParam(r, "productID"))
	productID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID do produto inválido", http.StatusBadRequest)
		return
	}
	product, err := h.service.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar o produto", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Erro ao encodar a resposta JSON: %v", err)
	}
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(chi.URLParam(r, "productID"))
	productID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID do produto inválido", http.StatusBadRequest)
		return
	}
	var productFromRequest domain.Produto
	if err := json.NewDecoder(r.Body).Decode(&productFromRequest); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}
	updatedProduct, err := h.service.UpdateProduct(r.Context(), productID, productFromRequest)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao atualizar o produto", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedProduct); err != nil {
		log.Printf("Erro ao encodar a resposta JSON: %v", err)
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(chi.URLParam(r, "productID"))
	productID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID do produto inválido", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteProduct(r.Context(), productID)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao deletar o produto", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// TransferStock realiza a transferência de estoque global para o estoque de um cliente.
func (h *ProductHandler) TransferStock(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "ID do produto inválido", http.StatusBadRequest)
		return
	}

	var req service.TransferStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.TransferStock(r.Context(), productID, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Transferência de estoque realizada com sucesso."}); err != nil {
		log.Printf("Erro ao codificar JSON na resposta de transferência: %v", err)
	}
}
