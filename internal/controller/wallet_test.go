package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fardinabir/digital-wallet-demo/internal/db"
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/fardinabir/digital-wallet-demo/internal/repository"
	"github.com/fardinabir/digital-wallet-demo/internal/service"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestWalletHandler_Create(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)
	repository := repository.NewWallet(dbInstance)
	service := service.NewWallet(repository)
	handler := NewWallet(service)

	tests := []struct {
		name       string
		createBody string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_create_user_wallet",
			createBody: `{"user_id":"test-user-001", "acnt_type":"user"}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"user_id":"test-user-001", "acnt_type":"user", "balance":0, "status":"active"}}`),
			},
		},
		{
			name:       "successful_create_provider_wallet",
			createBody: `{"user_id":"test-provider-001", "acnt_type":"provider"}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"user_id":"test-provider-001", "acnt_type":"provider", "balance":0, "status":"active"}}`),
			},
		},
		{
			name:       "missing_user_id",
			createBody: `{"acnt_type":"user"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "missing_acnt_type",
			createBody: `{"user_id":"test-user-002"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "invalid_acnt_type",
			createBody: `{"user_id":"test-user-003", "acnt_type":"invalid"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "invalid_request_body",
			createBody: `{"user_id":123}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean database before each test
			clearDB(dbInstance, model.Wallet{}, model.Transaction{})

			// Prepare
			req := httptest.NewRequest(http.MethodPost, "/wallets", bytes.NewReader([]byte(tt.createBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets")

			// Execute
			require.NoError(t, handler.Create(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"created_at": 1, "updated_at": 1, "id": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestWalletHandler_Deposit(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)
	repository := repository.NewWallet(dbInstance)
	service := service.NewWallet(repository)
	handler := NewWallet(service)

	tests := []struct {
		name        string
		setupWallet bool
		depositBody string
		want        want
	}{
		{
			name:        "successful_deposit",
			setupWallet: true,
			depositBody: `{"user_id":"test-user-001", "amount":5000, "provider_id":"deposit-provider-master"}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"subject_wallet_id":"test-user-001", "object_wallet_id":"deposit-provider-master", "transaction_type":"deposit", "operation_type":"credit", "amount":5000, "status":"completed"}}`),
			},
		},
		{
			name:        "deposit_without_provider_id",
			setupWallet: true,
			depositBody: `{"user_id":"test-user-001", "amount":3000}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"subject_wallet_id":"test-user-001", "object_wallet_id":"deposit-provider-master", "transaction_type":"deposit", "operation_type":"credit", "amount":3000, "status":"completed"}}`),
			},
		},
		{
			name:        "missing_user_id",
			setupWallet: false,
			depositBody: `{"amount":5000, "provider_id":"deposit-provider-master"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "missing_amount",
			setupWallet: true,
			depositBody: `{"user_id":"test-user-001", "provider_id":"deposit-provider-master"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "negative_amount",
			setupWallet: true,
			depositBody: `{"user_id":"test-user-001", "amount":-1000, "provider_id":"deposit-provider-master"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "zero_amount",
			setupWallet: true,
			depositBody: `{"user_id":"test-user-001", "amount":0, "provider_id":"deposit-provider-master"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "wallet_not_found",
			setupWallet: false,
			depositBody: `{"user_id":"non-existent-user", "amount":5000, "provider_id":"deposit-provider-master"}`,
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean database before each test
			clearDB(dbInstance, model.Wallet{}, model.Transaction{})

			// Setup wallet if needed
			if tt.setupWallet {
				createTestWallet(t, dbInstance, "test-user-001", model.User)
				createTestWalletWithBalance(t, dbInstance, "deposit-provider-master", model.Provider, 1000000) // Provider with sufficient balance
			}
			// Always create default provider for tests that might use it
			if tt.name == "deposit_without_provider_id" && !tt.setupWallet {
				createTestWallet(t, dbInstance, "test-user-001", model.User)
				createTestWalletWithBalance(t, dbInstance, "deposit-provider-master", model.Provider, 1000000)
			}

			// Prepare
			req := httptest.NewRequest(http.MethodPost, "/wallets/deposit", bytes.NewReader([]byte(tt.depositBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets/deposit")

			// Execute
			require.NoError(t, handler.Deposit(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"created_at": 1, "updated_at": 1, "id": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestWalletHandler_Withdraw(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)
	repository := repository.NewWallet(dbInstance)
	service := service.NewWallet(repository)
	handler := NewWallet(service)

	tests := []struct {
		name           string
		setupWallet    bool
		initialBalance int64
		withdrawBody   string
		want           want
	}{
		{
			name:           "successful_withdraw",
			setupWallet:    true,
			initialBalance: 10000,
			withdrawBody:   `{"user_id":"test-user-001", "amount":3000, "provider_id":"withdraw-provider-master"}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"subject_wallet_id":"test-user-001", "object_wallet_id":"withdraw-provider-master", "transaction_type":"withdraw", "operation_type":"debit", "amount":3000, "status":"completed"}}`),
			},
		},
		{
			name:           "insufficient_funds",
			setupWallet:    true,
			initialBalance: 1000,
			withdrawBody:   `{"user_id":"test-user-001", "amount":5000, "provider_id":"withdraw-provider-master"}`,
			want: want{
				StatusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:         "missing_user_id",
			setupWallet:  false,
			withdrawBody: `{"amount":3000, "provider_id":"withdraw-provider-master"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:           "missing_amount",
			setupWallet:    true,
			initialBalance: 10000,
			withdrawBody:   `{"user_id":"test-user-001", "provider_id":"withdraw-provider-master"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:           "negative_amount",
			setupWallet:    true,
			initialBalance: 10000,
			withdrawBody:   `{"user_id":"test-user-001", "amount":-1000, "provider_id":"withdraw-provider-master"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:         "wallet_not_found",
			setupWallet:  false,
			withdrawBody: `{"user_id":"non-existent-user", "amount":3000, "provider_id":"withdraw-provider-master"}`,
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean database before each test
			clearDB(dbInstance, model.Wallet{}, model.Transaction{})

			// Setup wallet if needed
			if tt.setupWallet {
				createTestWalletWithBalance(t, dbInstance, "test-user-001", model.User, tt.initialBalance)
				createTestWalletWithBalance(t, dbInstance, "withdraw-provider-master", model.Provider, 1000000) // Provider with sufficient balance
			}

			// Prepare
			req := httptest.NewRequest(http.MethodPost, "/wallets/withdraw", bytes.NewReader([]byte(tt.withdrawBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets/withdraw")

			// Execute
			require.NoError(t, handler.Withdraw(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"created_at": 1, "updated_at": 1, "id": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestWalletHandler_Transfer(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)
	repository := repository.NewWallet(dbInstance)
	service := service.NewWallet(repository)
	handler := NewWallet(service)

	tests := []struct {
		name         string
		setupWallets bool
		fromBalance  int64
		toBalance    int64
		transferBody string
		want         want
	}{
		{
			name:         "successful_transfer",
			setupWallets: true,
			fromBalance:  10000,
			toBalance:    5000,
			transferBody: `{"from_user_id":"test-user-001", "to_user_id":"test-user-002", "amount":3000}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"subject_wallet_id":"test-user-001", "object_wallet_id":"test-user-002", "transaction_type":"transfer", "operation_type":"debit", "amount":3000, "status":"completed"}}`),
			},
		},
		{
			name:         "insufficient_funds",
			setupWallets: true,
			fromBalance:  1000,
			toBalance:    5000,
			transferBody: `{"from_user_id":"test-user-001", "to_user_id":"test-user-002", "amount":5000}`,
			want: want{
				StatusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:         "transfer_to_same_wallet",
			setupWallets: true,
			fromBalance:  10000,
			toBalance:    0,
			transferBody: `{"from_user_id":"test-user-001", "to_user_id":"test-user-001", "amount":3000}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:         "missing_from_user_id",
			setupWallets: false,
			transferBody: `{"to_user_id":"test-user-002", "amount":3000}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:         "missing_to_user_id",
			setupWallets: false,
			transferBody: `{"from_user_id":"test-user-001", "amount":3000}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:         "missing_amount",
			setupWallets: true,
			fromBalance:  10000,
			toBalance:    5000,
			transferBody: `{"from_user_id":"test-user-001", "to_user_id":"test-user-002"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:         "negative_amount",
			setupWallets: true,
			fromBalance:  10000,
			toBalance:    5000,
			transferBody: `{"from_user_id":"test-user-001", "to_user_id":"test-user-002", "amount":-1000}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:         "from_wallet_not_found",
			setupWallets: false,
			transferBody: `{"from_user_id":"non-existent-user", "to_user_id":"test-user-002", "amount":3000}`,
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean database before each test
			clearDB(dbInstance, model.Wallet{}, model.Transaction{})

			// Setup wallets if needed
			if tt.setupWallets {
				createTestWalletWithBalance(t, dbInstance, "test-user-001", model.User, tt.fromBalance)
				if tt.name != "transfer_to_same_wallet" {
					createTestWalletWithBalance(t, dbInstance, "test-user-002", model.User, tt.toBalance)
				}
			}

			// Prepare
			req := httptest.NewRequest(http.MethodPost, "/wallets/transfer", bytes.NewReader([]byte(tt.transferBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets/transfer")

			// Execute
			require.NoError(t, handler.Transfer(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"created_at": 1, "updated_at": 1, "id": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestWalletHandler_Find(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	require.NoError(t, err)
	repository := repository.NewWallet(dbInstance)
	service := service.NewWallet(repository)
	handler := NewWallet(service)

	tests := []struct {
		name        string
		setupWallet bool
		userID      string
		want        want
	}{
		{
			name:        "successful_find",
			setupWallet: true,
			userID:      "test-user-001",
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":{"wallet":{"user_id":"test-user-001", "acnt_type":"user", "balance":10000, "status":"active"}, "transactions":[]}}`),
			},
		},
		{
			name:        "wallet_not_found",
			setupWallet: false,
			userID:      "non-existent-user",
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:        "missing_user_id",
			setupWallet: false,
			userID:      "",
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean database before each test
			clearDB(dbInstance, model.Wallet{}, model.Transaction{})

			// Setup wallet if needed
			if tt.setupWallet {
				createTestWalletWithBalance(t, dbInstance, "test-user-001", model.User, 10000)
			}

			// Prepare
			req := httptest.NewRequest(http.MethodGet, "/wallets/"+tt.userID, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets/:user_id")
			c.SetParamNames("user_id")
			c.SetParamValues(tt.userID)

			// Execute
			require.NoError(t, handler.Find(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"created_at": 1, "updated_at": 1, "id": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

// Helper functions
func clearDB(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model)
	}
}

func createTestWallet(t *testing.T, db *gorm.DB, userID string, acntType model.AcntType) {
	wallet := model.NewWallet(userID, acntType)
	err := db.Create(wallet).Error
	require.NoError(t, err)
}

func createTestWalletWithBalance(t *testing.T, db *gorm.DB, userID string, acntType model.AcntType, balance int64) {
	wallet := model.NewWallet(userID, acntType)
	wallet.Balance = balance
	err := db.Create(wallet).Error
	require.NoError(t, err)
}
