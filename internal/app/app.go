package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"thumbnail/internal/database/sqlite"
	"thumbnail/internal/repositories"
	"thumbnail/internal/services"
	transportGRPC "thumbnail/internal/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer *grpc.Server
	port       int
	db         *sqlite.Connection
}

func New() (*App, error) {
	port, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		return nil, fmt.Errorf("environment variable 'GRPC_PORT' must be type of INT and not nil: %w", err)
	}

	db := sqlite.NewConnection("thumbnails.db")
	repos := repositories.NewRepositories(db.DB)
	services := services.NewServices(repos)
	handlers := transportGRPC.NewHandler(services)

	grpcServer := grpc.NewServer()
	handlers.RegisterYoutubeHandler(grpcServer)
	reflection.Register(grpcServer)

	return &App{
		grpcServer: grpcServer,
		port:       port,
		db:         db,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Starting gRPC server on port %d", a.port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return a.grpcServer.Serve(lis)
}

func (a *App) Stop() {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}
	if a.db != nil {
		a.db.Close()
	}
}
