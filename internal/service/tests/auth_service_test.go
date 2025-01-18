package tests

import (
	"auth-service/internal/grpc"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq" // ou o driver do banco de dados que você está usando
)

func setupTestDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "user", "password", "auth_db",
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Importante: verifica se a conexão está funcionando
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}

func TestRegister(t *testing.T) {
	// Configurar o banco de dados
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	// Criar instância do repositório com a conexão do banco
	userRepo := repository.NewUserRepository(db)

	// Criar serviço com o repositório
	authService := service.NewAuthServiceServer(userRepo)

	req := &grpc.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Executar o teste
	resp, err := authService.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if resp.UserId == "" {
		t.Error("Expected user ID, got empty string")
	}
}
