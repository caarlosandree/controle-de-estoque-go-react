package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"controle-de-estoque/backend/internal/handler"
	"controle-de-estoque/backend/internal/repository"
	"controle-de-estoque/backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config representa a configuração da aplicação.
type Config struct {
	ServerAddress string
	DBURL         string
	CORSOrigins   []string
	JWTSecret     string
	Env           string
}

// Services agrupa todos os serviços da aplicação para fácil injeção.
type Services struct {
	TokenService   service.TokenGenerator
	UserService    *service.UserService
	ProductService *service.ProductService
	ClientService  *service.ClientService
}

// Handlers agrupa todos os handlers da aplicação.
type Handlers struct {
	ProductHandler *handler.ProductHandler
	UserHandler    *handler.UserHandler
	ClientHandler  *handler.ClientHandler
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Falha ao inicializar o logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()
	zap.ReplaceGlobals(logger)

	cfg, err := loadConfig()
	if err != nil {
		logger.Fatal("Falha ao carregar a configuração", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbpool, err := repository.NewDBConnection(ctx, cfg.DBURL)
	if err != nil {
		logger.Fatal("Falha ao conectar ao banco de dados",
			zap.Error(err),
			zap.String("dbURL", maskDBURL(cfg.DBURL)),
		)
	}
	defer dbpool.Close()
	logger.Info("Conexão com o banco de dados estabelecida")

	services := initServices(dbpool, cfg)
	handlers := initHandlers(services)

	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      setupRouter(handlers, services.TokenService, cfg),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	runServer(server, logger)
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v", err)
	}

	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL é obrigatório")
	}
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET é obrigatório")
	}
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DBURL:         dbURL,
		CORSOrigins:   strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
		JWTSecret:     jwtSecret,
		Env:           getEnv("ENV", "development"),
	}, nil
}

func initServices(dbpool *pgxpool.Pool, cfg *Config) *Services {
	productRepo := repository.NewProductRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)
	clientRepo := repository.NewClientRepository(dbpool)
	clientStockRepo := repository.NewClientStockRepository(dbpool) // ✅ corrigido para passar dbpool

	passwordService := service.NewPasswordService()
	tokenService := service.NewTokenService(cfg.JWTSecret)

	productService := service.NewProductService(dbpool, productRepo, clientStockRepo)
	userService := service.NewUserService(userRepo, passwordService, tokenService)
	clientService := service.NewClientService(clientRepo, clientStockRepo) // ✅ recebe estoque

	return &Services{
		TokenService:   tokenService,
		UserService:    userService,
		ProductService: productService,
		ClientService:  clientService,
	}
}

func initHandlers(s *Services) *Handlers {
	return &Handlers{
		ProductHandler: handler.NewProductHandler(s.ProductService),
		UserHandler:    handler.NewUserHandler(s.UserService, zap.L()),
		ClientHandler:  handler.NewClientHandler(s.ClientService),
	}
}

func setupRouter(h *Handlers, tokenService service.TokenGenerator, cfg *Config) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rotas públicas
	r.Get("/healthcheck", healthCheckHandler)
	r.Post("/register", h.UserHandler.Register)
	r.Post("/login", h.UserHandler.Login)

	// Rotas protegidas
	r.Group(func(r chi.Router) {
		r.Use(handler.AuthMiddleware(tokenService))

		r.Get("/me", h.UserHandler.GetMe)

		r.Route("/products", func(r chi.Router) {
			r.Post("/", h.ProductHandler.CreateProduct)
			r.Get("/", h.ProductHandler.ListProducts)
			r.Get("/{productID}", h.ProductHandler.GetProductByID)
			r.Put("/{productID}", h.ProductHandler.UpdateProduct)
			r.Delete("/{productID}", h.ProductHandler.DeleteProduct)
			r.Post("/{productID}/transfer", h.ProductHandler.TransferStock)
		})

		r.Route("/clients", func(r chi.Router) {
			r.Post("/", h.ClientHandler.CreateClient)
			r.Get("/", h.ClientHandler.ListClients)
			r.Get("/{clientID}", h.ClientHandler.GetClientByID)
			r.Put("/{clientID}", h.ClientHandler.UpdateClient)
			r.Delete("/{clientID}", h.ClientHandler.DeleteClient)

			// ✅ Nova rota de estoque do cliente
			r.Get("/{clientID}/stock", h.ClientHandler.ListStockByClientID)
		})
	})

	return r
}

func runServer(server *http.Server, logger *zap.Logger) {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Fatal("HTTP server shutdown error", zap.Error(err))
		}
		serverStopCtx()
	}()

	logger.Info("Starting HTTP server", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("HTTP server failed", zap.Error(err))
	}

	<-serverCtx.Done()
	logger.Info("Server stopped gracefully")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		zap.L().Error("failed to write health check response", zap.Error(err))
	}
}

// Helpers
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func maskDBURL(dbURL string) string {
	if strings.Contains(dbURL, "@") {
		parts := strings.Split(dbURL, "@")
		return "postgres://*****:*****@" + parts[1]
	}
	return dbURL
}
