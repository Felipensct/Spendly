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
	"github.com/stretchr/testify/assert"
)

// setupTestDB mantido como está pois está bem estruturado
func setupTestDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "user", "password", "auth_db",
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}

func setupService(t *testing.T) (*service.AuthServiceServer, func()) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthServiceServer(userRepo)

	cleanup := func() {
		// Limpa a tabela de usuários
		_, err := db.Exec("DELETE FROM users")
		if err != nil {
			t.Errorf("Failed to cleanup test database: %v", err)
		}
		db.Close()
	}

	return authService, cleanup
}

// TestRegister mantido como teste de caso feliz
func TestRegister(t *testing.T) {
	authService, cleanup := setupService(t)
	defer cleanup()

	req := &grpc.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := authService.Register(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.UserId)
}

func TestRegisterUsernameValidation(t *testing.T) {
	authService, cleanup := setupService(t)
	defer cleanup()

	testCases := []struct {
		name     string
		username string
		errMsg   string
	}{
		{
			name:     "empty username",
			username: "",
			errMsg:   "erro de username: o nome de usuário deve ter no mínimo 3 caracteres",
		},
		{
			name:     "username too short",
			username: "ab",
			errMsg:   "erro de username: o nome de usuário deve ter no mínimo 3 caracteres",
		},
		{
			name:     "username too long",
			username: "thisusernameiswaytoolongforoursystem123",
			errMsg:   "erro de username: o nome de usuário deve ter no máximo 30 caracteres",
		},
		{
			name:     "username with invalid characters",
			username: "user@name",
			errMsg:   "erro de username: username deve conter apenas letras, números, underscores e hifens",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &grpc.RegisterRequest{
				Username: tc.username,
				Email:    "valid@email.com",
				Password: "ValidPass123!",
			}

			resp, err := authService.Register(context.Background(), req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
			assert.Nil(t, resp)
		})
	}
}

func TestRegisterEmailValidation(t *testing.T) {
	authService, cleanup := setupService(t)
	defer cleanup()

	testCases := []struct {
		name   string
		email  string
		errMsg string
	}{
		{
			name:   "invalid email format",
			email:  "invalid-email",
			errMsg: "email inválido",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &grpc.RegisterRequest{
				Username: "validuser",
				Email:    tc.email,
				Password: "ValidPass123!",
			}

			resp, err := authService.Register(context.Background(), req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
			assert.Nil(t, resp)
		})
	}
}

func TestRegisterPasswordValidation(t *testing.T) {
	authService, cleanup := setupService(t)
	defer cleanup()

	testCases := []struct {
		name     string
		password string
		errMsg   string
	}{
		{
			name:     "password too short",
			password: "Short1!",
			errMsg:   "erro de senha: a senha deve ter no mínimo 8 caracteres",
		},
		{
			name:     "password without uppercase",
			password: "password123!",
			errMsg:   "erro de senha: a senha deve conter pelo menos uma letra maiúscula",
		},
		// ... outros casos de senha
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &grpc.RegisterRequest{
				Username: "validuser",
				Email:    "valid@email.com",
				Password: tc.password,
			}

			resp, err := authService.Register(context.Background(), req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
			assert.Nil(t, resp)
		})
	}
}

func TestRegisterSuccess(t *testing.T) {
	authService, cleanup := setupService(t)
	defer cleanup()

	req := &grpc.RegisterRequest{
		Username: "validuser",
		Email:    "valid@email.com",
		Password: "ValidPass123!",
	}

	resp, err := authService.Register(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.UserId)
	assert.Equal(t, "Usuário criado com sucesso", resp.Message)
}
