package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var DB *sql.DB

func InitDB(host, port, user, password, dbname string) {
	// Garantir que a porta não está vazia
	if port == "" {
		port = "5432" // porta padrão do PostgreSQL
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	// Tenta conectar ao banco com retry
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil {
			// Testa a conexão
			err = DB.Ping()
			if err == nil {
				log.Println("Successfully connected to database")
				return
			}
		}
		log.Printf("Failed to connect to database. Retrying in 2 seconds... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("Could not connect to the database: %v", err)
}
