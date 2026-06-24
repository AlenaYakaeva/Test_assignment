package server

import (
	"context"
	"exampleApp/internal/server/wallet"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv *http.Server
}

func New(addr string, walletService wallet.WalletService) *Server {
	srv := &http.Server{
		Addr: addr,
	}
	wh := wallet.NewWalletHandler(walletService)
	r := configureRouter(wh)
	srv.Handler = r

	return &Server{
		srv: srv,
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func configureRouter(wh *wallet.WalletHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/api/v1/wallet", wh.ChageAmount)
	r.GET("/api/v1/wallets/:id", wh.GetBalanceByID)

	return r
}
