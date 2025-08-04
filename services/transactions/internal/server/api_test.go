package server

import (
	"testing"

	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewAPI(t *testing.T) {
	tests := []struct {
		name    string
		opts    TxnAPIServerOpts
		wantErr bool
	}{
		{
			name: "Valid configuration",
			opts: TxnAPIServerOpts{
				ListenPort: 8082,
				Config: model.Config{
					PostgreSQL: model.PostgreSQL{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "postgres",
						DBName:   "wallet_test",
						SSLMode:  "disable",
					},
					SwaggerServer: model.Server{
						Enable: true,
						Port:   8082,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid database configuration",
			opts: TxnAPIServerOpts{
				ListenPort: 8082,
				Config: model.Config{
					PostgreSQL: model.PostgreSQL{
						Host:     "invalid-host",
						Port:     5432,
						User:     "invalid-user",
						Password: "invalid-password",
						DBName:   "invalid-db",
						SSLMode:  "disable",
					},
					SwaggerServer: model.Server{
						Enable: true,
						Port:   8082,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewAPI(tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, server)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, tt.opts.ListenPort, server.(*txnAPIServer).port)
				assert.IsType(t, &echo.Echo{}, server.(*txnAPIServer).engine)
				assert.IsType(t, &log.Entry{}, server.(*txnAPIServer).log)
				assert.IsType(t, &gorm.DB{}, server.(*txnAPIServer).db)
			}
		})
	}
}
