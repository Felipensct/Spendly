package repository

import (
	"auth-service/internal/utils"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &UserRepository{
		db: db,
	}
}

// CreateUser cria um novo usuário no banco de dados
func (repository *UserRepository) CreateUser(username, email, password string) (int, error) {
	// Validações adicionais de segurança
	if valid, err := utils.ValidateUsername(username); !valid {
		return 0, fmt.Errorf("username inválido: %s", err)
	}

	if !utils.ValidateEmail(email) {
		return 0, errors.New("email inválido")
	}
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err := repository.db.QueryRow(query, username, email, password).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// GetUserByUsername busca um usuário no banco de dados pelo email
func (repository *UserRepository) GetUserByUsername(username string) (int, string, error) {
	query := `SELECT id, password_hash FROM users WHERE username = $1`
	var userId int
	var password string
	err := repository.db.QueryRow(query, username).Scan(&userId, &password)
	if err != nil {
		return 0, "", err
	}
	return userId, password, nil
}
