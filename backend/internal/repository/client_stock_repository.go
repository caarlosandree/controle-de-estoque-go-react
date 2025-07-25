package repository

import (
	"context"
	"fmt"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool" // Importando o pgxpool
)

// ClientStockRepository gerencia as operações de banco de dados para o estoque de clientes.
type ClientStockRepository struct {
	db *pgxpool.Pool // Adicionando a dependência do banco
}

// NewClientStockRepository cria uma nova instância.
func NewClientStockRepository(db *pgxpool.Pool) *ClientStockRepository {
	return &ClientStockRepository{db: db}
}

// Upsert atualiza a quantidade de estoque de um cliente ou insere um novo registro.
func (r *ClientStockRepository) Upsert(ctx context.Context, tx pgx.Tx, stock *domain.ClientStock) error {
	query := `
		INSERT INTO client_stocks (client_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (client_id, product_id) DO UPDATE
		SET quantity = client_stocks.quantity + EXCLUDED.quantity
	` // Corrigido para somar a nova quantidade
	_, err := tx.Exec(ctx, query, stock.ClientID, stock.ProductID, stock.Quantity)
	if err != nil {
		return fmt.Errorf("erro ao fazer upsert no estoque do cliente: %w", err)
	}
	return nil
}

// ListStockByClientID busca o estoque de um cliente, juntando dados do produto.
func (r *ClientStockRepository) ListStockByClientID(ctx context.Context, clientID uuid.UUID) ([]domain.ClientStockDetails, error) {
	query := `
		SELECT
			cs.client_id,
			cs.product_id,
			p.name AS product_name,
			cs.quantity
		FROM
			client_stocks cs
		JOIN
			products p ON cs.product_id = p.id
		WHERE
			cs.client_id = $1
		ORDER BY
			p.name ASC
	`
	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar estoque do cliente: %w", err)
	}
	defer rows.Close()

	stocks := make([]domain.ClientStockDetails, 0)
	for rows.Next() {
		var s domain.ClientStockDetails
		if err := rows.Scan(&s.ClientID, &s.ProductID, &s.ProductName, &s.Quantity); err != nil {
			return nil, fmt.Errorf("erro ao escanear estoque do cliente: %w", err)
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}
