package server

import (
	"testing"

	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewAPI(t *testing.T) {
	tests := []struct {
		name    string
		opts    WalletAPIServerOpts
		wantErr bool
	}{
		{
			name: "Valid configuration",
			opts: WalletAPIServerOpts{
				ListenPort: 8081,
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
						Port:   8081,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid database configuration",
			opts: WalletAPIServerOpts{
				ListenPort: 8081,
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
						Port:   8081,
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
				assert.Equal(t, tt.opts.ListenPort, server.(*walletAPIServer).port)
				assert.IsType(t, &echo.Echo{}, server.(*walletAPIServer).engine)
				assert.IsType(t, &log.Entry{}, server.(*walletAPIServer).log)
				assert.IsType(t, &gorm.DB{}, server.(*walletAPIServer).db)
			}
		})
	}
}
