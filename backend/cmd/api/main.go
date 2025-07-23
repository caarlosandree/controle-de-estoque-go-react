package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"controle-de-estoque/backend/internal/handler"
	"controle-de-estoque/backend/internal/repository"
	"controle-de-estoque/backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors" // <<-- Importando o pacote CORS
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo .env.")
	}

	ctx := context.Background()

	dbpool, err := repository.NewDBConnection(ctx)
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer dbpool.Close()

	log.Println("Conexão com o banco de dados estabelecida com sucesso.")

	productRepo := repository.NewProductRepository(dbpool)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r := chi.NewRouter()

	// --- Configuração dos Middlewares ---

	// Middleware de CORS
	r.Use(cors.Handler(cors.Options{
		// Lista de origens permitidas. Nosso frontend está em http://localhost:5173
		AllowedOrigins: []string{"http://localhost:5173"},
		// Métodos HTTP que o frontend pode usar
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Cabeçalhos HTTP que podem ser enviados
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Tempo máximo que o resultado de uma pre-flight request pode ser cacheado
	}))

	// Middleware para logging de requisições
	r.Use(middleware.Logger)

	// --- Configuração das Rotas ---
	r.Get("/healthcheck", healthCheckHandler)

	r.Route("/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.ListProducts)
		r.Get("/{productID}", productHandler.GetProductByID)
		r.Put("/{productID}", productHandler.UpdateProduct)
		r.Delete("/{productID}", productHandler.DeleteProduct)
	})

	port := ":8080"
	log.Printf("Servidor da API iniciando na porta %s", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Não foi possível iniciar o servidor: %s\n", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API está funcionando normalmente e conectada ao banco de dados.")
}
