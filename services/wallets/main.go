// Package main is the entry point of the application.
package main

import (
	"github.com/fardinabir/digital-wallet-demo/services/wallets/cmd"
	_ "github.com/fardinabir/digital-wallet-demo/services/wallets/docs"
)

// @title			digital-wallet-demonstration API
// @version		0.0.1
// @description	This is a server for digital-wallet-demonstration.
// @license.name	Apache 2.0
// @host			localhost:8081
// @BasePath		/api/v1
// @schemes		http
func main() {
	cmd.Execute()
}
