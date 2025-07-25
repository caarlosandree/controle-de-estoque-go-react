package service

import (
	"context"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// IClientRepository define a interface para o repositório de clientes.
type IClientRepository interface {
	CreateClient(ctx context.Context, client *domain.Client) error
	ListClients(ctx context.Context) ([]domain.Client, error)
	GetClientByID(ctx context.Context, clientID uuid.UUID) (*domain.Client, error)
	UpdateClient(ctx context.Context, client *domain.Client) error
	DeleteClient(ctx context.Context, clientID uuid.UUID) error
}

// IClientStockRepository define a interface para o repositório de estoque do cliente.
type IClientStockRepository interface {
	ListStockByClientID(ctx context.Context, clientID uuid.UUID) ([]domain.ClientStockDetails, error)
	Upsert(ctx context.Context, tx pgx.Tx, stock *domain.ClientStock) error
}

// ClientService contém a lógica de negócio para clientes e estoques dos clientes.
type ClientService struct {
	repo      IClientRepository
	stockRepo IClientStockRepository
}

// NewClientService cria uma nova instância de ClientService.
func NewClientService(repo IClientRepository, stockRepo IClientStockRepository) *ClientService {
	return &ClientService{
		repo:      repo,
		stockRepo: stockRepo,
	}
}

// Create cria um novo cliente.
func (s *ClientService) Create(ctx context.Context, client *domain.Client) error {
	return s.repo.CreateClient(ctx, client)
}

// List retorna todos os clientes.
func (s *ClientService) List(ctx context.Context) ([]domain.Client, error) {
	return s.repo.ListClients(ctx)
}

// GetByID retorna um cliente pelo ID.
func (s *ClientService) GetByID(ctx context.Context, clientID uuid.UUID) (*domain.Client, error) {
	return s.repo.GetClientByID(ctx, clientID)
}

// Update atualiza os dados de um cliente.
func (s *ClientService) Update(ctx context.Context, client *domain.Client) error {
	return s.repo.UpdateClient(ctx, client)
}

// Delete remove um cliente pelo ID.
func (s *ClientService) Delete(ctx context.Context, clientID uuid.UUID) error {
	return s.repo.DeleteClient(ctx, clientID)
}

// ListStockByClientID retorna os dados de estoque de um cliente específico.
func (s *ClientService) ListStockByClientID(ctx context.Context, clientID uuid.UUID) ([]domain.ClientStockDetails, error) {
	return s.stockRepo.ListStockByClientID(ctx, clientID)
}
