package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"pvZ/init/migrations"
	grpcapi "pvZ/internal/adapters/api/grpc"
	"pvZ/internal/adapters/api/grpc/pvzpb"
	"pvZ/internal/config"
	"pvZ/internal/logger"
	"pvZ/internal/metrics"
)

func RunPVZ() {
	logger.Init()
	metrics.Init()

	cfg := config.Load()

	dbx := sqlx.MustConnect("postgres", cfg.DB.DSN())
	defer dbx.Close()

	if err := migrations.RunMigrations(dbx.DB); err != nil {
		logger.Log.Error("migration failed", "error", err)
	}

	secretKey := []byte(cfg.App.JWTSecret)
	deps := SetupDependencies(dbx, secretKey)

	// REST API
	go func() {
		addr := ":" + cfg.App.Rest
		router := SetupRoutes(deps.UserUC, deps.PVZUC, deps.ReceptionUC, deps.ProductUC, secretKey)
		logger.Log.Info("REST API server running on " + addr)
		if err := http.ListenAndServe(addr, router); err != nil {
			logger.Log.Error("REST API server error", "error", err)
		}
	}()

	// Prometheus
	go func() {
		addr := ":" + cfg.App.Prometheus
		http.Handle("/metrics", promhttp.Handler())
		logger.Log.Info("Prometheus metrics server running on " + addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			logger.Log.Error("Prometheus server error", "error", err)
		}
	}()

	// gRPC
	go func() {
		addr := ":" + cfg.App.GRPC
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Log.Error("gRPC listen error", "error", err)
			return
		}
		grpcServer := grpc.NewServer()
		pvzpb.RegisterPVZServiceServer(grpcServer, grpcapi.NewPVZService(deps.PVZUC))
		logger.Log.Info("gRPC server running on " + addr)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Log.Error("gRPC server error", "error", err)
		}
	}()

	select {}
}
