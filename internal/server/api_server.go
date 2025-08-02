package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// walletAPIServer is the API server for Wallet
type walletAPIServer struct {
	port   int
	engine *echo.Echo
	log    *log.Entry
	db     *gorm.DB
}

func (s *walletAPIServer) Name() string {
	return "walletAPIServer"
}

// Run starts the Wallet API server
func (s *walletAPIServer) Run() error {
	log.Infof("%s serving on port %d", s.Name(), s.port)
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

// Shutdown stops the Wallet API server
func (s *walletAPIServer) Shutdown(ctx context.Context) error {
	log.Infof("shutting down %s serving on port %d", s.Name(), s.port)
	return s.engine.Shutdown(ctx)
}
