package domain

import "errors"

// Erros comuns do domínio
var (
	ErrUserNotFound        = errors.New("usuário não encontrado")
	ErrEmailAlreadyExists  = errors.New("email já está em uso")
	ErrInvalidUserData     = errors.New("dados do usuário inválidos")
	ErrProductNotFound     = errors.New("produto não encontrado")
	ErrInvalidCredentials  = errors.New("credenciais inválidas")
	ErrUnauthorized        = errors.New("não autorizado")
	ErrInternalServerError = errors.New("erro interno do servidor")
)
