package repository

import (
	"context"
	"fmt"

	"controle-de-estoque/backend/internal/domain"

	"github.com/jackc/pgx/v5"
)

// ClientStockRepository gerencia as operações de banco de dados para o estoque de clientes.
type ClientStockRepository struct{}

// NewClientStockRepository cria uma nova instância.
func NewClientStockRepository() *ClientStockRepository {
	return &ClientStockRepository{}
}

// Upsert atualiza a quantidade de estoque de um cliente ou insere um novo registro.
// "Upsert" = UPDATE ou INSERT.
func (r *ClientStockRepository) Upsert(ctx context.Context, tx pgx.Tx, stock *domain.ClientStock) error {
	query := `
		INSERT INTO client_stocks (client_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (client_id, product_id) DO UPDATE
		SET quantity = client_stocks.quantity + $3
	`
	_, err := tx.Exec(ctx, query, stock.ClientID, stock.ProductID, stock.Quantity)
	if err != nil {
		return fmt.Errorf("erro ao fazer upsert no estoque do cliente: %w", err)
	}
	return nil
}
