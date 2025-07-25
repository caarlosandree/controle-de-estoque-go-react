package service

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService implementa operações seguras de hash e verificação de senhas
type PasswordService struct {
	minLength    int
	minUpperCase int
	minLowerCase int
	minDigits    int
	minSpecial   int
}

// NewPasswordService cria uma nova instância configurável
func NewPasswordService() *PasswordService {
	return &PasswordService{
		minLength:    8,
		minUpperCase: 1,
		minLowerCase: 1,
		minDigits:    1,
		minSpecial:   1,
	}
}

// Erros customizados de validação de senha
var (
	ErrPasswordTooShort      = errors.New("a senha deve ter no mínimo 8 caracteres")
	ErrPasswordTooWeak       = errors.New("a senha deve conter letras maiúsculas, minúsculas, números e caracteres especiais")
	ErrPasswordHashingFailed = errors.New("falha ao gerar hash da senha")
)

// HashPassword gera um hash seguro da senha usando bcrypt, após validar sua força
func (s *PasswordService) HashPassword(password string) (string, error) {
	if err := s.ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrPasswordHashingFailed
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash verifica se uma senha corresponde ao hash
func (s *PasswordService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength verifica os requisitos mínimos da senha
func (s *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < s.minLength {
		return ErrPasswordTooShort
	}

	var (
		hasUpper, hasLower, hasDigit, hasSpecial bool
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return ErrPasswordTooWeak
	}

	if s.isCommonPassword(password) {
		return ErrPasswordTooWeak
	}

	return nil
}

// isCommonPassword (exemplo básico)
func (s *PasswordService) isCommonPassword(password string) bool {
	commonPasswords := map[string]struct{}{
		"password": {}, "123456": {}, "qwerty": {}, "admin": {}, "welcome": {},
	}
	_, exists := commonPasswords[password]
	return exists
}
