// Package server provides the API server for the application.
package server

import (
	"fmt"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/controller"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/repository"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/service"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/utils"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/db"
	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// WalletAPIServerOpts is the options for the WalletAPIServer
type WalletAPIServerOpts struct {
	ListenPort int
	Config     model.Config
}

// NewAPI returns a new instance of the Wallet API server
func NewAPI(opts WalletAPIServerOpts) (Server, error) {
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

	s := &walletAPIServer{
		port:   opts.ListenPort,
		engine: engine,
		log:    logger,
		db:     dbInstance,
	}

	s.setupRoutes(engine)

	engine.Use(requestLogger())

	return s, nil
}

// initWalletController creates and configures the wallet handler with its dependencies
//
//	Repository ====> Service =====> Controller
//
// It follows the CSR dependency injection pattern
func (s *walletAPIServer) initWalletController() controller.WalletHandler {

	// Initialize dependencies (Repository -> Service -> Controller)
	walletRepo := repository.NewWalletRepo(s.db)
	walletService := service.NewWalletService(walletRepo)
	walletController := controller.NewWalletController(walletService)

	return walletController
}

// setupRoutes registers the routes for the application.
func (s *walletAPIServer) setupRoutes(e *echo.Echo) {
	e.Validator = controller.NewCustomValidator()

	api := e.Group("/api/v1")

	// Health check
	healthHandler := controller.NewHealth()
	api.GET("/health", healthHandler.Health)

	walletHandler := s.initWalletController()

	controller.InitRoutes(api, walletHandler)
}
