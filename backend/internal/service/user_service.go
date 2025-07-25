package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode"

	"controle-de-estoque/backend/internal/domain"

	"github.com/google/uuid"
)

// Interfaces de dependências
type (
	UserRepository interface {
		CreateUser(ctx context.Context, user *domain.User) error
		GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
		GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) // Corrigido: busca por ID
		UserExists(ctx context.Context, email string) (bool, error)
	}

	PasswordHasher interface {
		HashPassword(password string) (string, error)
		CheckPasswordHash(password, hash string) bool
		ValidatePasswordStrength(password string) error
	}

	TokenGenerator interface {
		GenerateToken(userID uuid.UUID) (string, error)
		ValidateToken(token string) (*uuid.UUID, error)
	}
)

// UserService implementa a lógica de negócio para usuários
type UserService struct {
	repo   UserRepository
	hasher PasswordHasher
	token  TokenGenerator
}

// NewUserService cria uma instância de UserService
func NewUserService(repo UserRepository, hasher PasswordHasher, token TokenGenerator) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
		token:  token,
	}
}

// DTOs (Data Transfer Objects)
type (
	RegisterRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	AuthResponse struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expiresAt"`
		UserID    uuid.UUID `json:"userId"`
	}
	UserProfile struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

// Erros customizados do serviço
var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrEmailInUse         = errors.New("email já está em uso")
	ErrWeakPassword       = errors.New("a senha não atende aos requisitos de segurança")
	ErrPasswordsDontMatch = errors.New("as senhas não coincidem")
)

// Register cria um novo usuário
func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	if req.Password != req.PasswordConfirm {
		return nil, ErrPasswordsDontMatch
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}
	exists, err := s.repo.UserExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar email: %w", err)
	}
	if exists {
		return nil, ErrEmailInUse
	}

	hashedPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	token, err := s.token.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	return &AuthResponse{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// Login autentica um usuário
func (s *UserService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if !s.hasher.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.token.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	return &AuthResponse{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// GetProfile retorna informações do perfil do usuário
func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	user, err := s.repo.GetUserByID(ctx, userID) // Corrigido para usar GetUserByID
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar perfil: %w", err)
	}

	return &UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// validatePassword verifica os requisitos de segurança da senha
func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	var (
		hasUpper, hasLower, hasDigit bool
	)
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return ErrWeakPassword
	}
	return nil
}
