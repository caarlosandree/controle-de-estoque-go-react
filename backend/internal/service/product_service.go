package service

import (
	"context"
	"controle-de-estoque/backend/internal/domain"
	"github.com/google/uuid"
	"math"
)

type ProductRepositoryInterface interface {
	CreateProduct(ctx context.Context, product *domain.Produto) error
	ListProducts(ctx context.Context, search string, page, limit int) ([]domain.Produto, int, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error)
	UpdateProduct(ctx context.Context, product *domain.Produto) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error // Novo método
}

type ProductService struct {
	repo ProductRepositoryInterface
}

func NewProductService(repo ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Produto) error {
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) ListProducts(ctx context.Context, search string, page, limit int) (*domain.PaginatedResponse, error) {
	products, totalRecords, err := s.repo.ListProducts(ctx, search, page, limit)
	if err != nil {
		return nil, err
	}

	// Calcula os metadados da paginação.
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

func (s *ProductService) GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error) {
	return s.repo.GetProductByID(ctx, productID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID uuid.UUID, input domain.Produto) (*domain.Produto, error) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	product.Name = input.Name
	product.Description = input.Description
	product.PriceInCents = input.PriceInCents
	product.Quantity = input.Quantity

	err = s.repo.UpdateProduct(ctx, &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// DeleteProduct chama o repositório para remover um produto.
func (s *ProductService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, productID)
}
