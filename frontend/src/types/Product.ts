// Esta interface deve espelhar a struct `domain.Produto` do nosso backend em Go.
export interface Product {
    id: string; // Em Go é uuid.UUID, mas em JSON/TS se torna uma string
    name: string;
    description: string;
    price_in_cents: number;
    quantity: number;
    created_at: string; // Em Go é time.Time, em JSON/TS vira uma string no formato ISO 8601
    updated_at: string;
}