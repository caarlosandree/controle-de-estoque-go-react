package repository

import (
	"context"
	"errors"
	"fmt"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClientRepository gerencia as operações de banco de dados para clientes.
type ClientRepository struct {
	db *pgxpool.Pool
}

// NewClientRepository cria uma nova instância de ClientRepository.
func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{db: db}
}

// CreateClient insere um novo cliente no banco de dados.
func (r *ClientRepository) CreateClient(ctx context.Context, client *domain.Client) error {
	query := `INSERT INTO clients (name, email, phone) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, client.Name, client.Email, client.Phone).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar cliente: %w", err)
	}
	return nil
}

// ListClients busca todos os clientes.
func (r *ClientRepository) ListClients(ctx context.Context) ([]domain.Client, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM clients ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar clientes: %w", err)
	}
	defer rows.Close()

	clients := make([]domain.Client, 0)
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("erro ao escanear cliente: %w", err)
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// GetClientByID busca um cliente pelo seu ID.
func (r *ClientRepository) GetClientByID(ctx context.Context, clientID uuid.UUID) (*domain.Client, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM clients WHERE id = $1`
	var c domain.Client
	err := r.db.QueryRow(ctx, query, clientID).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("cliente não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar cliente por ID: %w", err)
	}
	return &c, nil
}

// UpdateClient atualiza um cliente existente.
func (r *ClientRepository) UpdateClient(ctx context.Context, client *domain.Client) error {
	query := `UPDATE clients SET name = $1, email = $2, phone = $3, updated_at = NOW() WHERE id = $4 RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, client.Name, client.Email, client.Phone, client.ID).Scan(&client.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("cliente não encontrado para atualizar")
		}
		return fmt.Errorf("erro ao atualizar cliente: %w", err)
	}
	return nil
}

// DeleteClient remove um cliente do banco de dados.
func (r *ClientRepository) DeleteClient(ctx context.Context, clientID uuid.UUID) error {
	query := `DELETE FROM clients WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, clientID)
	if err != nil {
		return fmt.Errorf("erro ao deletar cliente: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("cliente não encontrado para deletar")
	}
	return nil
}
