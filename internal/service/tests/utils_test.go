package tests

import (
	"auth-service/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hashed := service.HashPassword(password)
	assert.NotEmpty(t, hashed, "Hashed password should not be empty")
}

func TestVerifyPassword(t *testing.T) {
	password := "securepassword"

	// Gera o hash da senha
	hashed := service.HashPassword(password)

	// Verifica se a senha correta passa
	if !service.VerifyPassword(hashed, password) {
		t.Fatalf("Falha na verificação: senha válida foi rejeitada")
	}

	// Verifica se uma senha incorreta falha
	if service.VerifyPassword(hashed, "wrongpassword") {
		t.Fatalf("Falha na verificação: senha inválida foi aceita")
	}
}
