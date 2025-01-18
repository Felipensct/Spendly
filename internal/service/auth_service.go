package service

import (
	"auth-service/internal/grpc"
	"auth-service/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//TODO: Estudar uso do regex para validação de email e senha

type AuthServiceServer struct {
	grpc.UnimplementedAuthServiceServer
	userRepo *repository.UserRepository
}

func NewAuthServiceServer(userRepo *repository.UserRepository) *AuthServiceServer {
	return &AuthServiceServer{
		userRepo: userRepo,
	}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// HashPassword gera um hash da senha do usuário
func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("erro ao gerar hash da senha: %v", err)
	}
	return string(hashed)
}

// VerifyPassword verifica se a senha do usuário é válida
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Register cria um novo usuário
func (s *AuthServiceServer) Register(ctx context.Context, req *grpc.RegisterRequest) (*grpc.RegisterResponse, error) {
	hashedPassword := HashPassword(req.Password)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("todos os campos são obrigatórios")
	}

	// Use s.userRepo ao invés da variável global userRepo
	userId, err := s.userRepo.CreateUser(req.Username, req.Email, hashedPassword)
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

	// Busca o usuário no banco de dados e verifica se a senha está correta
	userId, storedHashedPassword, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil || !VerifyPassword(storedHashedPassword, hashedPassword) {
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
