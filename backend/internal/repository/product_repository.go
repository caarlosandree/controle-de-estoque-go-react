package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5" // Import para ErrNoRows
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrProductNotFound é um erro customizado que retornamos quando um produto não é encontrado.
var ErrProductNotFound = errors.New("produto não encontrado")

// ProductRepository gerencia as operações de banco de dados para produtos.
type ProductRepository struct {
	db *pgxpool.Pool
}

// NewProductRepository cria uma nova instância de ProductRepository.
func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct insere um novo produto no banco de dados.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Produto) error {
	query := `
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

// ListProducts busca todos os produtos no banco de dados.
func (r *ProductRepository) ListProducts(ctx context.Context, search string, page, limit int) ([]domain.Produto, int, error) {
	// --- Primeira query: Contar o total de registros ---
	countQuery := "SELECT COUNT(*) FROM products"
	var countArgs []any
	if search != "" {
		countQuery += " WHERE name ILIKE $1" // ILIKE para busca case-insensitive
		countArgs = append(countArgs, "%"+search+"%")
	}

	var totalRecords int
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalRecords); err != nil {
		return nil, 0, fmt.Errorf("erro ao contar produtos: %w", err)
	}

	if totalRecords == 0 {
		return make([]domain.Produto, 0), 0, nil
	}

	// --- Segunda query: Buscar os dados da página atual ---
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id, name, description, price_in_cents, quantity, created_at, updated_at FROM products")

	var args []any
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

	products := make([]domain.Produto, 0)
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

// GetProductByID busca um único produto pelo seu ID.
func (r *ProductRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (domain.Produto, error) {
	query := `
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

// UpdateProduct atualiza um produto existente no banco de dados.
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Produto) error {
	query := `
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

// DeleteProduct remove um produto do banco de dados pelo seu ID.
func (r *ProductRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	query := "DELETE FROM products WHERE id = $1"
	commandTag, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}
