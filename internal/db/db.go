package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// Inicia a conexão com o banco de dados
func InitDB(host, port, user, password, dbname string) {
	// String de conexão com o POSTGRESQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error

	// Retry logic to ensure the database container is ready
	for i := 0; i < 10; i++ {
		DB, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database. Retrying in 2 seconds... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Database connected successfully!")
}
