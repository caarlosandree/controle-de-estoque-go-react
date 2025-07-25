package handler

// contextKey é um tipo privado para evitar colisões de chave no contexto.
type contextKey string

// UserIDContextKey é a chave usada para armazenar o ID do usuário no contexto da requisição.
const UserIDContextKey contextKey = "userID"
