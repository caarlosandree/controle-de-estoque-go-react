package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBConnection(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	if databaseUrl == "" {
		return nil, fmt.Errorf("a URL do banco de dados (DATABASE_URL) não foi definida")
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
