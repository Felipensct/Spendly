package repository

import (
	"auth-service/internal/db"
)

type UserRepository struct{}

// CreateUser cria um novo usuário no banco de dados
func (repository *UserRepository) CreateUser(username, email, password string) (int, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err := db.DB.QueryRow(query, username, email, password).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// GetUserByUsername busca um usuário no banco de dados pelo email
func (repository *UserRepository) GetUserByUsername(username string) (int, string, error) {
	query := `SELECT id, password FROM users WHERE name = $1`
	var userId int
	var password string
	err := db.DB.QueryRow(query, username).Scan(&userId, &password)
	if err != nil {
		return 0, "", err
	}
	return userId, password, nil
}
