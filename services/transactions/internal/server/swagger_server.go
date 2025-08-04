package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// swaggerServer is the API server for Txn
type swaggerServer struct {
	port   int
	engine *echo.Echo
	log    *log.Entry
}

// SwaggerServerOpts is the options for the swaggerServer
type SwaggerServerOpts struct {
	ListenPort int
}

// NewSwagger returns a new instance of the Swagger server
func NewSwagger(opts SwaggerServerOpts) Server {
	logger := log.NewEntry(log.StandardLogger())
	log.SetFormatter(&log.JSONFormatter{})

	engine := echo.New()
	engine.HideBanner = true

	engine.Use(requestLogger())

	engine.GET("/swagger/*", echoSwagger.WrapHandler)

	s := &swaggerServer{
		port:   opts.ListenPort,
		engine: engine,
		log:    logger,
	}

	return s
}

func (s *swaggerServer) Name() string {
	return "swaggerServer"
}

func (s *swaggerServer) Run() error {
	log.Infof("%s serving on port %d", s.Name(), s.port)
	return s.engine.Start(fmt.Sprintf(":%d", s.port))
}

func (s *swaggerServer) Shutdown(ctx context.Context) error {
	log.Infof("shutting down %s serving on port %d", s.Name(), s.port)
	return s.engine.Shutdown(ctx)
}
