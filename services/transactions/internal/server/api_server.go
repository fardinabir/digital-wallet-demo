package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// walletAPIServer is the API server for Txn
type txnAPIServer struct {
	port   int
	engine *echo.Echo
	log    *log.Entry
	db     *gorm.DB
}

func (s *txnAPIServer) Name() string {
	return "txnAPIServer"
}

// Run starts the Txn API server
func (s *txnAPIServer) Run() error {
	log.Infof("%s serving on port %d", s.Name(), s.port)
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

// Shutdown stops the Txn API server
func (s *txnAPIServer) Shutdown(ctx context.Context) error {
	log.Infof("shutting down %s serving on port %d", s.Name(), s.port)
	return s.engine.Shutdown(ctx)
}
