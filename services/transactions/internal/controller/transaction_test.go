package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/db"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/model"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/repository"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/service"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTransactionHandler_CreateTransactionPair(t *testing.T) {
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
	repository := repository.NewTransactionRepository(dbInstance)
	service := service.NewTransactionService(repository)
	handler := NewTransactionHandler(service)

	tests := []struct {
		name       string
		createBody string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_create_transaction_pair",
			createBody: `{"debit_transaction":{"subject_wallet_id":"user-001","object_wallet_id":"user-002","transaction_type":"transfer","operation_type":"debit","amount":1000,"status":"completed"},"credit_transaction":{"subject_wallet_id":"user-002","object_wallet_id":"user-001","transaction_type":"transfer","operation_type":"credit","amount":1000,"status":"completed"}}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":"Transaction pair created successfully"}`),
			},
		},
		{
			name:       "successful_deposit_transaction_pair",
			createBody: `{"debit_transaction":{"subject_wallet_id":"deposit-provider-master","object_wallet_id":"user-001","transaction_type":"deposit","operation_type":"debit","amount":5000,"status":"completed"},"credit_transaction":{"subject_wallet_id":"user-001","object_wallet_id":"deposit-provider-master","transaction_type":"deposit","operation_type":"credit","amount":5000,"status":"completed"}}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":"Transaction pair created successfully"}`),
			},
		},
		{
			name:       "missing_debit_transaction",
			createBody: `{"credit_transaction":{"subject_wallet_id":"user-002","object_wallet_id":"user-001","transaction_type":"transfer","operation_type":"credit","amount":1000,"status":"completed"}}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "missing_credit_transaction",
			createBody: `{"debit_transaction":{"subject_wallet_id":"user-001","object_wallet_id":"user-002","transaction_type":"transfer","operation_type":"debit","amount":1000,"status":"completed"}}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "invalid_amount_zero",
			createBody: `{"debit_transaction":{"subject_wallet_id":"user-001","object_wallet_id":"user-002","transaction_type":"transfer","operation_type":"debit","amount":0,"status":"completed"},"credit_transaction":{"subject_wallet_id":"user-002","object_wallet_id":"user-001","transaction_type":"transfer","operation_type":"credit","amount":0,"status":"completed"}}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "invalid_request_body",
			createBody: `{"invalid":"json"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean database before each test
			clearDB(dbInstance, model.Transaction{})

			// Prepare
			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(tt.createBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/transactions")

			// Execute
			require.NoError(t, handler.CreateTransactionPair(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestTransactionHandler_GetTransactions(t *testing.T) {
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
	repository := repository.NewTransactionRepository(dbInstance)
	service := service.NewTransactionService(repository)
	handler := NewTransactionHandler(service)

	tests := []struct {
		name             string
		setupTransaction bool
		subjectWalletID  string
		want             want
	}{
		{
			name:             "successful_get_transactions",
			setupTransaction: true,
			subjectWalletID:  "user-001",
			want: want{
				StatusCode: http.StatusOK,
			},
		},
		{
			name:             "get_transactions_empty_result",
			setupTransaction: false,
			subjectWalletID:  "user-002",
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":[]}`),
			},
		},
		{
			name:            "missing_subject_wallet_id",
			subjectWalletID: "",
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean database before each test
			clearDB(dbInstance, model.Transaction{})

			// Setup transaction if needed
			if tt.setupTransaction {
				createTestTransaction(t, dbInstance, "user-001", "user-002", model.Transfer, model.Debit, 1000)
			}

			// Prepare
			req := httptest.NewRequest(http.MethodGet, "/transactions/"+tt.subjectWalletID, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/transactions/:subject_wallet_id")
			c.SetParamNames("subject_wallet_id")
			c.SetParamValues(tt.subjectWalletID)

			// Execute
			require.NoError(t, handler.GetTransactions(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
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

func createTestTransaction(t *testing.T, db *gorm.DB, subjectWalletID, objectWalletID string, txnType model.TransactionType, opType model.OperationType, amount int64) {
	txn := model.NewTransaction(subjectWalletID, objectWalletID, txnType, opType, amount)
	txn.Status = model.Completed
	err := db.Create(txn).Error
	require.NoError(t, err)
}
