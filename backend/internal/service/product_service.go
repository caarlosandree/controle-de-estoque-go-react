package service

import (
	"context"
	"math"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
)

// ProductRepositoryInterface define o contrato que nosso serviço espera do repositório.
type ProductRepositoryInterface interface {
	CreateProduct(ctx context.Context, product *domain.Produto) error
	ListProducts(ctx context.Context, search string, page, limit int) ([]domain.Produto, int, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error)
	UpdateProduct(ctx context.Context, product *domain.Produto) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
}

// ProductService contém a lógica de negócio para produtos.
type ProductService struct {
	repo ProductRepositoryInterface
}

// NewProductService cria uma nova instância de ProductService.
func NewProductService(repo ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct chama o repositório para criar um produto.
func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Produto) error {
	return s.repo.CreateProduct(ctx, product)
}

// ListProducts busca os produtos e calcula os metadados de paginação.
func (s *ProductService) ListProducts(ctx context.Context, search string, page, limit int) (*domain.PaginatedResponse, error) {
	products, totalRecords, err := s.repo.ListProducts(ctx, search, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if totalRecords > 0 {
		totalPages = int(math.Ceil(float64(totalRecords) / float64(limit)))
	}

	metadata := domain.Metadata{
		TotalRecords: totalRecords,
		CurrentPage:  page,
		PageSize:     limit,
		TotalPages:   totalPages,
	}

	response := &domain.PaginatedResponse{
		Data:     products,
		Metadata: metadata,
	}

	return response, nil
}

// GetProductByID chama o repositório para buscar um produto por ID.
func (s *ProductService) GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error) {
	return s.repo.GetProductByID(ctx, productID)
}

// UpdateProduct busca um produto, atualiza seus campos e o salva.
func (s *ProductService) UpdateProduct(ctx context.Context, productID uuid.UUID, input domain.Produto) (*domain.Produto, error) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	product.Name = input.Name
	product.Description = input.Description
	product.PriceInCents = input.PriceInCents
	product.Quantity = input.Quantity

	if err := s.repo.UpdateProduct(ctx, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

// DeleteProduct chama o repositório para remover um produto.
func (s *ProductService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, productID)
}
