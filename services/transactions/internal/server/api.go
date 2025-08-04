// Package server provides the API server for the application.
package server

import (
	"fmt"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/controller"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/repository"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/service"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/utils"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/db"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// TxnAPIServerOpts is the options for the TxnAPIServer
type TxnAPIServerOpts struct {
	ListenPort int
	Config     model.Config
}

// NewAPI returns a new instance of the Txn API server
func NewAPI(opts TxnAPIServerOpts) (Server, error) {
	logger := log.NewEntry(log.StandardLogger())
	log.SetFormatter(&log.JSONFormatter{})

	// Initialize global logger
	utils.InitLogger(logger)

	dbInstance, err := db.New(opts.Config.PostgreSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	engine := echo.New()

	// Allow all origins for CORS
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	s := &txnAPIServer{
		port:   opts.ListenPort,
		engine: engine,
		log:    logger,
		db:     dbInstance,
	}

	s.setupRoutes(engine)

	engine.Use(requestLogger())

	return s, nil
}

// initTransactionController creates and configures the transaction handler with its dependencies
//
//	Repository ====> Service =====> Controller
//
// It follows the CSR dependency injection pattern
func (s *txnAPIServer) initTransactionController() controller.TransactionHandler {

	// Initialize dependencies (Repository -> Service -> Controller)
	transactionRepo := repository.NewTransactionRepository(s.db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionController := controller.NewTransactionHandler(transactionService)

	return transactionController
}

// setupRoutes registers the routes for the application.
func (s *txnAPIServer) setupRoutes(e *echo.Echo) {
	e.Validator = controller.NewCustomValidator()

	api := e.Group("/api/v1")

	// Health check
	healthHandler := controller.NewHealth()
	api.GET("/health", healthHandler.Health)

	transactionHandler := s.initTransactionController()

	controller.InitRoutes(api, transactionHandler)
}
