package handler

import (
	"context"
	"net/http"
	"strings"

	"controle-de-estoque/backend/internal/service"
)

// AuthMiddleware é um middleware para proteger rotas.
func AuthMiddleware(tokenService *service.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Cabeçalho de autorização ausente", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Cabeçalho de autorização mal formatado", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Valida o token e recupera o userID
			userIDPtr, err := tokenService.ValidateToken(tokenString)
			if err != nil || userIDPtr == nil {
				http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
				return
			}

			// Insere userID (string) no contexto da requisição
			ctx := context.WithValue(r.Context(), UserIDContextKey, userIDPtr.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
