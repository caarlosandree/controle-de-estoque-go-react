package domain

import (
	"time"

	"github.com/google/uuid"
)

// ClientStock representa a tabela de ligação.
type ClientStock struct {
	ClientID  uuid.UUID `json:"client_id" db:"client_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ClientStockDetails é um DTO para a resposta da API, incluindo o nome do produto.
type ClientStockDetails struct {
	ClientID    uuid.UUID `json:"clientId"`
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
}
