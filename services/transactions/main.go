package main

import (
	"github.com/fardinabir/digital-wallet-demo/services/transactions/cmd"
	_ "github.com/fardinabir/digital-wallet-demo/services/transactions/docs"
)

// @title			Digital Txn Transactions API
// @version			v1.0
// @description		This is a digital wallet transactions microservice API.
// @license.name		MIT
// @license.url		https://opensource.org/licenses/MIT
// @host			localhost:8082
// @BasePath		/api/v1
// @schemes			http https
func main() {
	cmd.Execute()
}
