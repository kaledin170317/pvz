package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"pvZ/internal/metrics"

	"pvZ/init/migrations"
	grpcapi "pvZ/internal/adapters/api/grpc"
	"pvZ/internal/adapters/api/grpc/pvzpb"
	"pvZ/internal/config"
)

func RunPVZ() {
	cfg := config.Load()

	dbx := sqlx.MustConnect("postgres", cfg.DB.DSN())
	defer dbx.Close()

	if err := migrations.RunMigrations(dbx.DB); err != nil {
		log.Println("migration error:", err)
	}

	secretKey := []byte(cfg.App.JWTSecret)
	deps := SetupDependencies(dbx, secretKey)

	metrics.Init()

	// REST API
	go func() {
		router := SetupRoutes(deps.UserUC, deps.PVZUC, deps.ReceptionUC, deps.ProductUC, secretKey)
		log.Println("REST API server running on :8080")
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	// Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics server running on :9000")
		log.Fatal(http.ListenAndServe(":9000", nil))
	}()

	// gRPC
	go func() {
		lis, err := net.Listen("tcp", ":3000")
		if err != nil {
			log.Fatalf("gRPC listen error: %v", err)
		}
		grpcServer := grpc.NewServer()
		pvzpb.RegisterPVZServiceServer(grpcServer, grpcapi.NewPVZService(deps.PVZUC))
		log.Println("gRPC server running on :3000")
		log.Fatal(grpcServer.Serve(lis))
	}()

	select {}
}
