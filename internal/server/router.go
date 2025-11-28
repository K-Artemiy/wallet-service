package server

import (
	"context"
	"net/http"
	"time"
	"wallet-service/internal/config"
	"wallet-service/internal/repository/pgrepo"
	"wallet-service/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	serv *http.Server
}

func NewServer(pool *pgxpool.Pool) *Server {
	mux := http.NewServeMux()

	repo := pgrepo.NewWalletRepository(pool)
	uc := usecase.NewWalletUseCase(repo)
	h := NewHandlers(uc)

	mux.HandleFunc("POST /api/v1/wallets", h.CreateOperationHandler)
	mux.HandleFunc("GET /api/v1/wallets/{walletId}", h.GetBalanceHandler)

	return &Server{
		serv: &http.Server{
			Addr:         config.GetAppPort(),
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.serv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.serv.Shutdown(ctx)
}
