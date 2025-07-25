package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository define a interface para operações com usuários
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UserExists(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository cria uma nova instância do UserRepository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// CreateUser cria um novo usuário no banco de dados
func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Violação de constraint única
				return domain.ErrEmailAlreadyExists
			case "23514": // Violação de check constraint
				return domain.ErrInvalidUserData
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail busca um usuário pelo email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetUserByID busca um usuário pelo ID
func (r *userRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// UpdateUser atualiza os dados de um usuário
func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, updated_at = $3 
		WHERE id = $4
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// DeleteUser remove um usuário do banco de dados
func (r *userRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// UserExists verifica se um usuário com o email especificado existe
func (r *userRepository) UserExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}
