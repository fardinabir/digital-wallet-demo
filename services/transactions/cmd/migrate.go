// Package cmd provides the command line interface for the application.
package cmd

import (
	"fmt"

	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Run: func(_ *cobra.Command, _ []string) {
		dbInstance, err := db.New(cfg.PostgreSQL)
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err)
			return
		}

		if err := db.Migrate(dbInstance); err != nil {
			log.Fatalf("failed to migrate database: %s", err)
		}
		fmt.Printf("Migration completed. PostgreSQL database: %s\n", cfg.PostgreSQL.DBName)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
