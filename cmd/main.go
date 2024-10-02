package main

import (
	"context"
	"fmt"
	"garantex/internal"
	"garantex/internal/handler"
	"garantex/internal/pb"
	"garantex/internal/repository"
	"garantex/internal/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	_ "go.uber.org/automaxprocs"
)

const (
	serviceName = "garantex-proxy"
)

func main() {
	config := new(internal.Config)
	config.Read()
	config.Print()

	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	db, err := pgxpool.New(context.Background(), config.Conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	r := repository.New(db)
	s := service.New(r)
	h := handler.New(s, logger)

	metrics := prometheus.NewServerMetrics()

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
	)

	//healthcheck
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)
	//prometheus metrics
	metrics.InitializeMetrics(grpcServer)
	http.Handle("/metrics", promhttp.Handler())
	//setting handler
	pb.RegisterCryptoExchangeServiceServer(grpcServer, h)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCPort))
	if err != nil {
		logger.Fatal("Listening gRPC error", zap.Error(err))
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterCryptoExchangeServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.GRPCPort), opts)
	if err != nil {
		logger.Fatal("Registering gRPC gateway endpoint error", zap.Error(err))
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	go func() {
		logger.Sugar().Infof("GRPC server is listening on :%d", config.GRPCPort)
		err := grpcServer.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
			log.Fatal(err)
		}
	}()

	go func() {
		logger.Sugar().Infof("GRPC gateway server is listening on :%d", config.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	go func() {
		logger.Sugar().Infof("Metrics server is listening on :%d", config.MetricsPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.MetricsPort), nil); err != nil {
			log.Fatal(err)
		}
	}()

	<-shutdown

	logger.Info("Shutdown signal received")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("GRPC gateway server shutdown error")
	}
	grpcServer.GracefulStop()
	logger.Info("Server stopped gracefully")
}
