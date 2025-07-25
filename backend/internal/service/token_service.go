package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenService é responsável por operações com tokens JWT
type TokenService struct {
	secretKey     []byte
	issuer        string
	tokenDuration time.Duration
}

// Claims personalizadas que estendem RegisteredClaims padrão
type customClaims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}

// Configuração padrão
const (
	defaultIssuer     = "controle-de-estoque-api"
	defaultDuration   = 15 * time.Minute   // Access token de curta duração
	refreshDuration   = 7 * 24 * time.Hour // Refresh token de longa duração
	tokenCookieName   = "access_token"
	refreshCookieName = "refresh_token"
)

// Erros customizados
var (
	ErrInvalidToken            = errors.New("token inválido")
	ErrUnexpectedSigningMethod = errors.New("método de assinatura inesperado")
	ErrTokenExpired            = errors.New("token expirado")
)

// NewTokenService cria uma nova instância de TokenService
func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secretKey:     []byte(secret),
		issuer:        defaultIssuer,
		tokenDuration: defaultDuration,
	}
}

// TokenPair representa um par de tokens (access + refresh)
type TokenPair struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// GenerateTokenPair gera um par de tokens (access + refresh)
func (s *TokenService) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	accessToken, expiresAt, err := s.generateToken(userID, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar access token: %w", err)
	}

	refreshToken, _, err := s.generateToken(userID, refreshDuration)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// generateToken gera um token JWT individual
func (s *TokenService) generateToken(userID uuid.UUID, duration time.Duration) (string, time.Time, error) {
	expirationTime := time.Now().Add(duration)
	claims := &customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

// ValidateToken verifica e decodifica um token JWT
func (s *TokenService) ValidateToken(tokenString string) (*uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return &claims.UserID, nil
	}

	return nil, ErrInvalidToken
}

// SetTokenCookies define os cookies de autenticação na resposta HTTP
func (s *TokenService) SetTokenCookies(w http.ResponseWriter, tokenPair *TokenPair) {
	accessCookie := &http.Cookie{
		Name:     tokenCookieName,
		Value:    tokenPair.AccessToken,
		Expires:  tokenPair.ExpiresAt,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Lax é melhor para SPA
	}

	refreshCookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    tokenPair.RefreshToken,
		Expires:  time.Now().Add(refreshDuration),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		Path:     "/api/refresh", // Escopo mais restrito para o refresh token
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}

// GenerateToken (para compatibilidade com a interface antiga, se necessário)
func (s *TokenService) GenerateToken(userID uuid.UUID) (string, error) {
	token, _, err := s.generateToken(userID, s.tokenDuration)
	return token, err
}
