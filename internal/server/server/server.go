package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"tradeservice/internal/config"
)

type Server struct {
	server *http.Server
	logger *slog.Logger
	port   string
}

func New(logger *slog.Logger, cfg *config.ServerConfig) *Server {

	server := &http.Server{
		Addr: ":" + cfg.Port,
		//Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	handler.GetRequest(w, r)
		//}),
	}

	return &Server{
		logger: logger,
		server: server,
		port:   cfg.Port,
	}
}
func (s Server) Run() {
	s.logger.Info("Server is running on: localhost", "Port", s.port)
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server starting error: %v", err)
		}
	}
}

func (s Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server...")
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error("Error: ", err)
	}

	return err
}
