package main

import (
	"auth-service/internal/db"
	pb "auth-service/internal/grpc"
	"auth-service/internal/service"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq" // Importa o driver PostgreSQL

	"google.golang.org/grpc"
)

func main() {
	//Inicializa a conexão com o banco
	db.InitDB(
		os.Getenv("DB_HOST"),     // "postgres" em vez de "localhost"
		os.Getenv("DB_PORT"),     // porta do banco
		os.Getenv("DB_USER"),     // usuário do banco
		os.Getenv("DB_PASSWORD"), // senha do banco
		os.Getenv("DB_NAME"),     // nome do banco
	)
	log.Println("Conexão com o banco inicializada")

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao escutar a porta 5432: %v", err)
	}

	grpcServer := grpc.NewServer()
	log.Println("gRPC server inicializado")

	pb.RegisterAuthServiceServer(grpcServer, &service.AuthServiceServer{})
	log.Println("AuthService registrado com o gRPC server")

	// Inicia o servidor gRPC
	log.Println("AuthService está sendo executado na porta 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
