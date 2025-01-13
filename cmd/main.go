package main

import (
	"auth-service/internal/db"
	pb "auth-service/internal/grpc"
	"auth-service/internal/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	//Inicializa a conexão com o banco
	db.InitDB("postgres", "5432", "user", "password", "auth_db")
	log.Println("Conexão com o banco inicializada")

	listener, err := net.Listen("tcp", "50051")
	if err != nil {
		log.Fatalf("Falha ao escutar a porta 50051: %v", err)
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
