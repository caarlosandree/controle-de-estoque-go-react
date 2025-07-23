package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBConnection cria e retorna um pool de conexões com o banco de dados PostgreSQL.
// Sua única responsabilidade é estabelecer a conexão.
func NewDBConnection(ctx context.Context) (*pgxpool.Pool, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, fmt.Errorf("a variável de ambiente DATABASE_URL não foi definida")
	}

	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar o pool de conexões com o banco: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("não foi possível conectar ao banco de dados: %w", err)
	}

	return pool, nil
}
