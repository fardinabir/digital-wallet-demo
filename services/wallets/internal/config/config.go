package config

import "github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"

var globalConfig *model.Config

// SetGlobalConfig sets the global configuration instance
func SetGlobalConfig(cfg *model.Config) {
	globalConfig = cfg
}

// GetGlobalConfig returns the global configuration instance
func GetGlobalConfig() *model.Config {
	return globalConfig
}
