package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"tradeservice/internal/config"
	"tradeservice/internal/server/handler"
	"tradeservice/internal/server/middleware"
	"tradeservice/internal/storage/postgres"
)

type Server struct {
	server  *echo.Echo
	logger  *slog.Logger
	storage *postgres.Storage
	port    int
}

func New(logger *slog.Logger,
	cfg *config.ServerConfig,
	db *postgres.Storage,
	categoryHandler *handlers.CategoriesController,
	productHandler *handlers.ProductController) *Server {

	server := echo.New()

	server.Use(middleware.LogRequest(logger))

	categoryGroup := server.Group("categories")

	categoryGroup.GET("", categoryHandler.Get)
	categoryGroup.DELETE("/:categoryId", categoryHandler.Delete)
	categoryGroup.POST("/create/:categoryName/:productId", categoryHandler.Add)
	categoryGroup.POST("/update/:categoryId/:categoryName", categoryHandler.Set)

	productGroup := server.Group("product")

	productGroup.GET("", productHandler.Get)
	productGroup.DELETE("/:productId", productHandler.Delete)
	productGroup.POST("/create/:productName/:productId", productHandler.Add)
	productGroup.POST("/update/:productName/:productId", productHandler.Set)

	return &Server{
		logger:  logger,
		server:  server,
		storage: db,
		port:    cfg.Port,
	}
}
func (s Server) Run() {
	s.logger.Info("Server is running on: localhost", "Port", s.port)
	if err := s.server.Start(fmt.Sprintf("localhost:%d", s.port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server starting error: %v", err)
		}
	}
}

func (s Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping DB Connection")

	s.storage.Close()

	s.logger.Info("Stopping server...")
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error("Error: ", err)
		return fmt.Errorf("error while stopping Server Request %w", err)
	}

	return nil
}
