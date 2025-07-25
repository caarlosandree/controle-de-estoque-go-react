package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IProductRepository define os métodos que o repositório de produtos deve implementar,
// incluindo os métodos para uso dentro de transação.
type IProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Produto) error
	ListProducts(ctx context.Context, search string, page, limit int) ([]domain.Produto, int, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error)
	UpdateProduct(ctx context.Context, product *domain.Produto) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error

	// Métodos para transação
	GetProductForUpdate(ctx context.Context, tx pgx.Tx, productID uuid.UUID) (*domain.Produto, error)
	UpdateQuantity(ctx context.Context, tx pgx.Tx, productID uuid.UUID, newQuantity int) error
}

// IClientStockRepository define o contrato para repositório que gerencia estoque por cliente.
type IClientStockRepository interface {
	Upsert(ctx context.Context, tx pgx.Tx, stock *domain.ClientStock) error
}

// ProductService contém a lógica de negócio para produtos, incluindo transferências de estoque.
type ProductService struct {
	db        *pgxpool.Pool // Pool para iniciar transações
	repo      IProductRepository
	stockRepo IClientStockRepository
}

// NewProductService cria uma instância de ProductService com as dependências necessárias.
func NewProductService(db *pgxpool.Pool, repo IProductRepository, stockRepo IClientStockRepository) *ProductService {
	return &ProductService{
		db:        db,
		repo:      repo,
		stockRepo: stockRepo,
	}
}

// CreateProduct cria um novo produto chamando o repositório.
func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Produto) error {
	return s.repo.CreateProduct(ctx, product)
}

// ListProducts busca produtos e retorna a resposta paginada.
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

	return &domain.PaginatedResponse{
		Data:     products,
		Metadata: metadata,
	}, nil
}

// GetProductByID busca um produto pelo ID.
func (s *ProductService) GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error) {
	return s.repo.GetProductByID(ctx, productID)
}

// UpdateProduct atualiza um produto.
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

// DeleteProduct remove um produto pelo ID.
func (s *ProductService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, productID)
}

// TransferStockRequest representa os dados para transferência de estoque a um cliente.
type TransferStockRequest struct {
	ClientID uuid.UUID `json:"clientId"`
	Quantity int       `json:"quantity"`
}

// TransferStock realiza a transferência de estoque global para o estoque de um cliente,
// garantindo atomicidade e consistência via transação.
func (s *ProductService) TransferStock(ctx context.Context, productID uuid.UUID, req TransferStockRequest) error {
	if req.Quantity <= 0 {
		return errors.New("a quantidade a ser transferida deve ser positiva")
	}

	// Inicia a transação
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // rollback silencioso caso não tenha commit
	}()

	// 1. Bloqueia o produto para update na transação
	product, err := s.repo.GetProductForUpdate(ctx, tx, productID)
	if err != nil {
		return err
	}

	// 2. Verifica estoque disponível
	if product.Quantity < req.Quantity {
		return fmt.Errorf("estoque insuficiente: disponível %d, solicitado %d", product.Quantity, req.Quantity)
	}

	// 3. Atualiza estoque global
	newQuantity := product.Quantity - req.Quantity
	if err := s.repo.UpdateQuantity(ctx, tx, productID, newQuantity); err != nil {
		return err
	}

	// 4. Atualiza estoque do cliente (upsert)
	clientStock := &domain.ClientStock{
		ClientID:  req.ClientID,
		ProductID: productID,
		Quantity:  req.Quantity,
	}
	if err := s.stockRepo.Upsert(ctx, tx, clientStock); err != nil {
		return err
	}

	// 5. Commit da transação
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}
