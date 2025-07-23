package domain

import (
	"time"

	"github.com/google/uuid"
)

// Produto representa a entidade de produto no nosso sistema.
// As tags `json` controlam como os campos são nomeados quando convertidos para JSON.
// As tags `db` serão usadas futuramente pela camada do banco de dados para mapear colunas.
type Produto struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	PriceInCents int64     `json:"price_in_cents" db:"price_in_cents"`
	Quantity     int       `json:"quantity" db:"quantity"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Explicação das Escolhas:
// - ID (uuid.UUID): Usar UUID como chave primária é uma ótima prática. Evita a adivinhação de IDs sequenciais
//   e facilita a vida em sistemas distribuídos. Precisaremos adicionar essa dependência.
//
// - PriceInCents (int64): NUNCA use float para dinheiro devido a problemas de arredondamento. A melhor prática
//   é armazenar o valor na menor unidade monetária (centavos) como um número inteiro.
//
// - CreatedAt / UpdatedAt (time.Time): Campos essenciais para auditoria. Sabemos quando um registro foi
//   criado e modificado pela última vez.
