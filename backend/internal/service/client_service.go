package service

import (
	"context"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
)

// IClientRepository define a interface para o repositório de clientes.
type IClientRepository interface {
	CreateClient(ctx context.Context, client *domain.Client) error
	ListClients(ctx context.Context) ([]domain.Client, error)
	GetClientByID(ctx context.Context, clientID uuid.UUID) (*domain.Client, error)
	UpdateClient(ctx context.Context, client *domain.Client) error
	DeleteClient(ctx context.Context, clientID uuid.UUID) error
}

// ClientService contém a lógica de negócio para clientes.
type ClientService struct {
	repo IClientRepository
}

// NewClientService cria uma nova instância de ClientService.
func NewClientService(repo IClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) Create(ctx context.Context, client *domain.Client) error {
	return s.repo.CreateClient(ctx, client)
}

func (s *ClientService) List(ctx context.Context) ([]domain.Client, error) {
	return s.repo.ListClients(ctx)
}

func (s *ClientService) GetByID(ctx context.Context, clientID uuid.UUID) (*domain.Client, error) {
	return s.repo.GetClientByID(ctx, clientID)
}

func (s *ClientService) Update(ctx context.Context, client *domain.Client) error {
	return s.repo.UpdateClient(ctx, client)
}

func (s *ClientService) Delete(ctx context.Context, clientID uuid.UUID) error {
	return s.repo.DeleteClient(ctx, clientID)
}
