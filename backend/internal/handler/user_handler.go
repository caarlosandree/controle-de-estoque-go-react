package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"controle-de-estoque/backend/internal/domain"
	"controle-de-estoque/backend/internal/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserHandler gerencia as requisições HTTP relacionadas a usuários
type UserHandler struct {
	userService *service.UserService
	logger      *zap.Logger
}

// NewUserHandler cria uma nova instância de UserHandler
func NewUserHandler(userService *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger.Named("UserHandler"),
	}
}

// RegisterRequest define a estrutura esperada para registro
type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

// LoginRequest define a estrutura esperada para login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ErrorResponse representa uma resposta de erro padrão
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// GetMe lida com a busca do perfil do usuário logado
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Pega o userID que o middleware colocou no contexto.
	userIDStr, ok := r.Context().Value(UserIDContextKey).(string)
	if !ok {
		h.sendError(w, "internal_error", "ID de usuário ausente no contexto", http.StatusInternalServerError)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.sendError(w, "invalid_user_id", "ID de usuário inválido", http.StatusBadRequest)
		return
	}

	userProfile, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			h.sendError(w, "user_not_found", "Usuário não encontrado", http.StatusNotFound)
			return
		}
		h.handleServiceError(w, err)
		return
	}

	h.sendJSON(w, userProfile, http.StatusOK)
}

// Register lida com requisições de registro de usuários
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invalid_request_body", "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	serviceReq := service.RegisterRequest{
		Email:           req.Email,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
	}

	authResponse, err := h.userService.Register(r.Context(), serviceReq)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSON(w, authResponse, http.StatusCreated)
}

// Login lida com requisições de autenticação de usuários
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invalid_request_body", "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	serviceReq := service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	authResponse, err := h.userService.Login(r.Context(), serviceReq)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSON(w, authResponse, http.StatusOK)
}

// handleServiceError trata erros retornados pelo serviço
func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		h.sendError(w, "invalid_credentials", "Credenciais inválidas", http.StatusUnauthorized)
	case errors.Is(err, service.ErrEmailInUse):
		h.sendError(w, "email_in_use", "Email já está em uso", http.StatusConflict)
	case errors.Is(err, service.ErrWeakPassword):
		h.sendError(w, "weak_password", err.Error(), http.StatusBadRequest)
	case errors.Is(err, service.ErrPasswordsDontMatch):
		h.sendError(w, "passwords_dont_match", "As senhas não coincidem", http.StatusBadRequest)
	default:
		h.logger.Error("Erro interno não mapeado", zap.Error(err))
		h.sendError(w, "internal_error", "Erro interno no servidor", http.StatusInternalServerError)
	}
}

// sendJSON envia uma resposta JSON
func (h *UserHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.logger.Error("Erro ao codificar resposta JSON", zap.Error(err))
		}
	}
}

// sendError envia uma resposta de erro padronizada
func (h *UserHandler) sendError(w http.ResponseWriter, errorCode, message string, statusCode int) {
	resp := ErrorResponse{
		Error:   errorCode,
		Message: message,
	}
	h.sendJSON(w, resp, statusCode)
}
