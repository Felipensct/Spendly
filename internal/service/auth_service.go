package service

import (
	"auth-service/internal/db"
	"auth-service/internal/grpc"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthServiceServer struct {
	grpc.UnimplementedAuthServiceServer
}

// TODO: Criar método mais segura de geração de chave
var jwtSecret = []byte("jwt_secret")

// HashPassword gera um hash da senha do usuário
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Register cria um novo usuário
func (s *AuthServiceServer) Register(ctx context.Context, req *grpc.RegisterRequest) (*grpc.RegisterResponse, error) {
	//Encripta a senha do usuário
	hashedPassword := HashPassword(req.Password)

	//TODO: Verificar se query dentro do serviço é uma boa prática
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err := db.DB.QueryRow(query, req.Username, req.Email, hashedPassword).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao inserir usuário: %v", err)
	}

	return &grpc.RegisterResponse{
		UserId:  fmt.Sprintf("%d", userId),
		Message: "Usuário criado com sucesso",
	}, nil
}

// Login realiza a autenticação do usuário
func (s *AuthServiceServer) Login(ctx context.Context, req *grpc.LoginRequest) (*grpc.LoginResponse, error) {
	hashedPassword := HashPassword(req.Password)
	var userId int

	query := `SELECT id, password_hash FROM users WHERE username = $1`
	err := db.DB.QueryRow(query, req.Username).Scan(&userId, &hashedPassword)
	if err != nil || hashedPassword != HashPassword(req.Password) {
		return nil, errors.New("usuário ou senha inválidos")
	}

	// Lógica para expiração do token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %v", err)
	}

	return &grpc.LoginResponse{
		Token: tokenString,
	}, nil
}

// ValidateToken valida o token do usuário
func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *grpc.ValidateTokenRequest) (*grpc.ValidateTokenResponse, error) {
	// Tenta analisar e validar o token JWT usando jwt.Parse com um callback personalizado
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return &grpc.ValidateTokenResponse{IsValid: false}, nil
	}

	// Extrai as claims do token JWT e verifica se o token é válido. Ok é true se o token é válido
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &grpc.ValidateTokenResponse{IsValid: false}, nil
	}

	return &grpc.ValidateTokenResponse{
		IsValid: true,
		UserId:  fmt.Sprintf("%v", claims["userId"]),
	}, nil
}
