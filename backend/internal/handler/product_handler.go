package handler

import (
	"controle-de-estoque/backend/internal/domain"
	"controle-de-estoque/backend/internal/repository"
	"controle-de-estoque/backend/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Lendo o parâmetro 'search'
	search := r.URL.Query().Get("search")

	// Lendo e validando o parâmetro 'page'
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // Valor padrão
	}

	// Lendo e validando o parâmetro 'limit'
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Valor padrão
	}

	response, err := h.service.ListProducts(r.Context(), search, page, limit)
	if err != nil {
		http.Error(w, "Erro ao listar os produtos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
	json.NewEncoder(w).Encode(product)
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
	json.NewEncoder(w).Encode(updatedProduct)
}

// DeleteProduct lida com a remoção de um produto.
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

	// Em um DELETE bem-sucedido, a resposta padrão é 204 No Content.
	w.WriteHeader(http.StatusNoContent)
}
