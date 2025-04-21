package app

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	"pvZ/init/migrations"
	grpcapi "pvZ/internal/adapters/api/grpc"
	"pvZ/internal/adapters/api/grpc/pvzpb"
	"pvZ/internal/config"
)

// Точка входа — запуск и REST, и gRPC
func RunPVZ() {
	cfg := config.Load()

	fmt.Println("DSN:", cfg.DB.DSN())
	dbx := sqlx.MustConnect("postgres", cfg.DB.DSN())
	defer dbx.Close()

	db := dbx.DB
	if err := migrations.RunMigrations(db); err != nil {
		log.Println("migration error:", err)
	}

	secretKey := []byte(cfg.App.JWTSecret)
	deps := SetupDependencies(dbx, secretKey)

	// REST-сервер
	go func() {
		router := SetupRoutes(deps.UserUC, deps.PVZUC, deps.ReceptionUC, deps.ProductUC, secretKey)
		fmt.Printf("🌐 REST-сервер на порту :%s\n", cfg.App.Port)
		log.Fatal(http.ListenAndServe(":"+cfg.App.Port, router))
	}()

	// gRPC-сервер
	go func() {
		lis, err := net.Listen("tcp", ":3000")
		if err != nil {
			log.Fatalf("не удалось слушать порт 3000: %v", err)
		}

		grpcServer := grpc.NewServer()
		pvzpb.RegisterPVZServiceServer(grpcServer, grpcapi.NewPVZService(deps.PVZUC))

		fmt.Println("🔌 gRPC сервер запущен на порту :3000")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC сервер завершился с ошибкой: %v", err)
		}
	}()

	select {}
}
