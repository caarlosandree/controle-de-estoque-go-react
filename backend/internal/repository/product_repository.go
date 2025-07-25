package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrProductNotFound é retornado quando um produto não é encontrado no banco.
var ErrProductNotFound = errors.New("produto não encontrado")

// ProductRepository gerencia operações no banco relacionadas a produtos.
type ProductRepository struct {
	db *pgxpool.Pool
}

// NewProductRepository cria um novo repositório de produtos.
func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetProductForUpdate busca um produto por ID e bloqueia a linha para update dentro da transação.
func (r *ProductRepository) GetProductForUpdate(ctx context.Context, tx pgx.Tx, productID uuid.UUID) (*domain.Produto, error) {
	const query = `
		SELECT id, name, description, price_in_cents, quantity
		FROM products
		WHERE id = $1
		FOR UPDATE
	`
	var p domain.Produto
	err := tx.QueryRow(ctx, query, productID).Scan(&p.ID, &p.Name, &p.Description, &p.PriceInCents, &p.Quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("erro ao buscar produto para atualização: %w", err)
	}
	return &p, nil
}

// UpdateQuantity atualiza a quantidade de um produto dentro de uma transação.
func (r *ProductRepository) UpdateQuantity(ctx context.Context, tx pgx.Tx, productID uuid.UUID, newQuantity int) error {
	const query = `
		UPDATE products
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id
	`
	var id uuid.UUID
	err := tx.QueryRow(ctx, query, newQuantity, productID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return fmt.Errorf("erro ao atualizar quantidade do produto: %w", err)
	}
	return nil
}

// CreateProduct insere um novo produto no banco.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Produto) error {
	const query = `
        INSERT INTO products (name, description, price_in_cents, quantity)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.PriceInCents,
		product.Quantity,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("não foi possível criar o produto: %w", err)
	}
	return nil
}

// ListProducts busca produtos com paginação e busca por nome (case-insensitive).
func (r *ProductRepository) ListProducts(ctx context.Context, search string, page, limit int) ([]domain.Produto, int, error) {
	countQuery := "SELECT COUNT(*) FROM products"
	var countArgs []any
	if search != "" {
		countQuery += " WHERE name ILIKE $1"
		countArgs = append(countArgs, "%"+search+"%")
	}

	var totalRecords int
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalRecords); err != nil {
		return nil, 0, fmt.Errorf("erro ao contar produtos: %w", err)
	}

	if totalRecords == 0 {
		return []domain.Produto{}, 0, nil
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT id, name, description, price_in_cents, quantity, created_at, updated_at
		FROM products
	`)

	args := []any{}
	argID := 1
	if search != "" {
		queryBuilder.WriteString(fmt.Sprintf(" WHERE name ILIKE $%d", argID))
		args = append(args, "%"+search+"%")
		argID++
	}

	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argID, argID+1))
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar produtos: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Produto, 0, limit)
	for rows.Next() {
		var p domain.Produto
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PriceInCents, &p.Quantity, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("erro ao escanear produto: %w", err)
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar pelos produtos: %w", err)
	}

	return products, totalRecords, nil
}

// GetProductByID busca um produto pelo ID.
func (r *ProductRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error) {
	const query = `
        SELECT id, name, description, price_in_cents, quantity, created_at, updated_at
        FROM products
        WHERE id = $1
    `
	var p domain.Produto
	err := r.db.QueryRow(ctx, query, productID).Scan(&p.ID, &p.Name, &p.Description, &p.PriceInCents, &p.Quantity, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Produto{}, ErrProductNotFound
		}
		return domain.Produto{}, fmt.Errorf("erro ao buscar produto por ID: %w", err)
	}
	return p, nil
}

// UpdateProduct atualiza os dados de um produto.
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Produto) error {
	const query = `
        UPDATE products
        SET name = $1, description = $2, price_in_cents = $3, quantity = $4, updated_at = NOW()
        WHERE id = $5
        RETURNING updated_at
    `
	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.PriceInCents,
		product.Quantity,
		product.ID,
	).Scan(&product.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return fmt.Errorf("erro ao atualizar produto: %w", err)
	}
	return nil
}

// DeleteProduct remove um produto pelo ID.
func (r *ProductRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	const query = "DELETE FROM products WHERE id = $1"
	commandTag, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}
