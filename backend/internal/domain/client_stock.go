package domain

import (
	"time"

	"github.com/google/uuid"
)

// ClientStock representa a quantidade de um produto específico para um cliente.
type ClientStock struct {
	ClientID  uuid.UUID `json:"client_id" db:"client_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
