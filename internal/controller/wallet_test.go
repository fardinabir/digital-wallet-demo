package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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
			name:       "successful_create",
			createBody: `{"task":"Created Task", "priority":1}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"Task":"Created Task", "Priority":1, "Status":"created"}}`),
			},
		},
		{
			name:       "successful_create_without_priority",
			createBody: `{"task":"Created Task"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "create_with_invalid_priority",
			createBody: `{"task":"Created Task", "priority":-1}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "successful_create_but_with_ignore_status",
			createBody: `{"task":"Created Task", "priority":2, "status":"done"}`,
			want: want{
				StatusCode: http.StatusCreated,
				Response:   []byte(`{"data":{"Task":"Created Task", "Priority":2, "Status":"created"}}`),
			},
		},
		{
			name:       "invalid_request_body",
			createBody: `{"task":1}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			req := httptest.NewRequest(http.MethodPost, "/dummy/target", bytes.NewReader([]byte(tt.createBody)))
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
				ignoreMapEntires(map[string]any{"CreatedAt": 1, "UpdatedAt": 1, "ID": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestWalletHandler_Update(t *testing.T) {
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
		updateBody string
		updateID   string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_update",
			createBody: `{"task":"Updated Task", "priority":1}`,
			updateBody: `{"task":"Updated Task","status":"done"}`,
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":{"Task":"Updated Task","Status":"done", "Priority":1}}`),
			},
		},
		{
			name:       "successful_update_with_priority",
			createBody: `{"task":"Updated Task", "priority":2}`,
			updateBody: `{"task":"Updated Task","status":"done", "priority":3}`,
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":{"Task":"Updated Task","Status":"done", "Priority":3}}`),
			},
		},
		{
			name:       "successful_update_with_only_priority",
			createBody: `{"task":"Updated Task", "priority":1}`,
			updateBody: `{"priority":2}`,
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":{"Task":"Updated Task","Status":"created", "Priority":2}}`),
			},
		},
		{
			name:       "update_with_invalid_priority",
			createBody: `{"task":"Updated Task", "priority":1}`,
			updateBody: `{"priority":22}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "update_with_invalid_status",
			createBody: `{"task":"Updated Task", "priority":1}`,
			updateBody: `{"task":"Updated Task","status":"pending"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "not_found_record",
			updateID:   "-1",
			updateBody: `{"task":"Updated Task","status":"done"}`,
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:       "invalid_request_body",
			updateBody: `{"task":1}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "invalid_request_parameter",
			updateID:   "invalid",
			updateBody: `{"task":"Updated Task","status":"done"}`,
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			id := "1"
			if tt.updateID != "" {
				id = tt.updateID
			} else if tt.createBody != "" {
				id = strconv.Itoa(createTask(t, e, handler, tt.createBody))
			}

			req := httptest.NewRequest(http.MethodPut, "/dummy/target", bytes.NewReader([]byte(tt.updateBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets/:id")
			c.SetParamNames("id")
			c.SetParamValues(id)

			// Execute
			require.NoError(t, handler.Update(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"CreatedAt": 1, "UpdatedAt": 1, "ID": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestWalletHandler_Delete(t *testing.T) {
	type want struct {
		StatusCode int
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
		deleteID   string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_delete",
			createBody: `{"task":"Deleted Task", "priority":1}`,
			want: want{
				StatusCode: http.StatusNoContent,
			},
		},
		{
			name:     "not_found_record",
			deleteID: "-1",
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:     "invalid_request_parameter",
			deleteID: "invalid",
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			id := "1"
			if tt.deleteID != "" {
				id = tt.deleteID
			} else if tt.createBody != "" {
				id = strconv.Itoa(createTask(t, e, handler, tt.createBody))
			}

			req := httptest.NewRequest(http.MethodDelete, "/dummy/target", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets/:id")
			c.SetParamNames("id")
			c.SetParamValues(id)

			// Execute
			require.NoError(t, handler.Delete(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)
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
		name       string
		createBody string
		findID     string
		want       want
		wantErr    bool
	}{
		{
			name:       "successful_find",
			createBody: `{"task":"Found Task", "priority":1}`,
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":{"Task":"Found Task","Status":"created", "Priority":1}}`),
			},
		},
		{
			name:   "not_found",
			findID: "-1",
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:   "invalid_request_parameter",
			findID: "invalid",
			want: want{
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare
			id := "1"
			if tt.findID != "" {
				id = tt.findID
			} else if tt.createBody != "" {
				id = strconv.Itoa(createTask(t, e, handler, tt.createBody))
			}

			req := httptest.NewRequest(http.MethodGet, "/dummy/target", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets/:id")
			c.SetParamNames("id")
			c.SetParamValues(id)

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
				ignoreMapEntires(map[string]any{"CreatedAt": 1, "UpdatedAt": 1, "ID": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func TestWalletHandler_FindAll(t *testing.T) {
	type want struct {
		StatusCode int
		Response   []byte
	}

	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	err = db.Migrate(dbInstance)
	clearDB(dbInstance, model.Wallet{})
	require.NoError(t, err)
	repository := repository.NewWallet(dbInstance)
	service := service.NewWallet(repository)
	handler := NewWallet(service)

	tests := []struct {
		name       string
		createBody []string
		want       want
	}{
		{
			name: "no_wallets_found",
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":[]}`), // Expecting an empty array
			},
		},
		{
			name:       "single_wallets_found",
			createBody: []string{`{"task":"Task 1", "priority":1}`},
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":[{"Task":"Task 1","Priority":1,"Status":"created"}]}`),
			},
		},
		{
			name:       "multiple_wallets_found",
			createBody: []string{`{"task":"Task 2", "priority":2}`},
			want: want{
				StatusCode: http.StatusOK,
				Response:   []byte(`{"data":[{"Task":"Task 2","Priority":2,"Status":"created"},{"Task":"Task 1","Priority":1,"Status":"created"}]}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a task if needed for the current test case
			for _, ct := range tt.createBody {
				createTask(t, e, handler, ct)
			}

			req := httptest.NewRequest(http.MethodGet, "/wallets", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/wallets")

			// Execute
			require.NoError(t, handler.FindAll(c))

			// Assert
			assert.Equal(t, tt.want.StatusCode, rec.Code)

			if tt.want.Response == nil {
				return
			}
			got := rec.Body.Bytes()

			opts := []cmp.Option{
				cmpTransformJSON(t),
				ignoreMapEntires(map[string]any{"CreatedAt": 1, "UpdatedAt": 1, "ID": 1}),
			}
			if diff := cmp.Diff(got, tt.want.Response, opts...); diff != "" {
				t.Errorf("return value mismatch (-got +want):\n%s", diff)
				t.Logf("got:\n%s", string(got))
			}
		})
	}
}

func clearDB(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model)
	}
}

func createTask(t *testing.T, e *echo.Echo, handler WalletHandler, body string) int {
	req := httptest.NewRequest(http.MethodPost, "/wallets", bytes.NewReader([]byte(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Create(c)
	require.NoError(t, err)

	type response struct {
		Data   model.Wallet
		Status string
	}

	var res response
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err, "Failed to json.Unmarshal err: %s", err)
	require.NotEmpty(t, res.Data.ID)

	return res.Data.ID
}
