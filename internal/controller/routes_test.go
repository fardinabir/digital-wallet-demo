package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fardinabir/digital-wallet-demo/internal/db"
	"github.com/fardinabir/digital-wallet-demo/internal/repository"
	"github.com/fardinabir/digital-wallet-demo/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)
	setupTestRoutes(e, dbInstance)

	// Test cases
	tests := []struct {
		name         string
		method       string
		target       string
		expectedCode int
	}{
		{"Health_Check", http.MethodGet, "/api/v1/health", http.StatusOK},
		{"Create_Wallet_without_body", http.MethodPost, "/api/v1/wallets", http.StatusBadRequest}, // Assuming no body is sent, should return BadRequest
		{"Get_all_Wallets", http.MethodGet, "/api/v1/wallets", http.StatusOK},
		{"Get_non-existent_Wallet", http.MethodGet, "/api/v1/wallets/1", http.StatusNotFound},       // Assuming no wallet with id 1 exists
		{"Update_Wallet_without_body", http.MethodPut, "/api/v1/wallets/1", http.StatusNotFound},    // Assuming no body is sent, should return BadRequest
		{"Delete_non-existent_Wallet", http.MethodDelete, "/api/v1/wallets/1", http.StatusNotFound}, // Assuming no wallet with id 1 exists
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

// setupTestRoutes configures routes for testing with the same pattern as the server
func setupTestRoutes(e *echo.Echo, db *gorm.DB) {
	// Set up request validation
	e.Validator = NewCustomValidator()

	// Create API version group
	api := e.Group("/api/v1")

	// Register health check endpoint
	healthHandler := NewHealth()
	api.GET("/health", healthHandler.Health)

	// Initialize wallet handler with dependencies
	walletRepo := repository.NewWallet(db)
	walletService := service.NewWallet(walletRepo)
	walletHandler := NewWallet(walletService)

	// Register wallet routes
	InitRoutes(api, walletHandler)
}
