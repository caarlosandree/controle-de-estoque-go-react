package domain

import (
	"time"

	"github.com/google/uuid"
)

// User representa a entidade de usuário no nosso sistema.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // O traço `json:"-"` impede que este campo seja exposto em respostas JSON.
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Explicação das Escolhas:
// - PasswordHash (string): Armazenaremos apenas o "hash" da senha, nunca a senha em texto plano.
//   Esta é a prática de segurança mais importante em sistemas de autenticação.
//   Usaremos uma biblioteca robusta para gerar este hash a partir da senha do usuário.
