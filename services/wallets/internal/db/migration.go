package db

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fardinabir/digital-wallet-demo/services/wallets/internal/model"
	"gorm.io/gorm"
)

// Migrate runs the complete migration process for the database
// It performs DDL migrations (schema) followed by DML migrations (data)
func Migrate(db *gorm.DB) error {
	// Step 1: Run GORM auto-migration for schema creation
	if err := db.AutoMigrate(&model.Wallet{}); err != nil {
		fmt.Printf("ERROR: Auto-migration failed: %v\n", err)
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}
	fmt.Println("Successfully completed GORM auto-migration")

	// Step 2: Run DDL migrations (if any)
	if err := runSQLMigrations(db, "migrations/ddl"); err != nil {
		fmt.Printf("ERROR: DDL migrations failed: %v\n", err)
		return fmt.Errorf("failed to run DDL migrations: %w", err)
	}
	fmt.Println("Successfully completed DDL migrations")

	// Step 3: Run DML migrations (data insertions)
	if err := runSQLMigrations(db, "migrations/dml"); err != nil {
		fmt.Printf("ERROR: DML migrations failed: %v\n", err)
		return fmt.Errorf("failed to run DML migrations: %w", err)
	}
	fmt.Println("Successfully completed DML migrations")

	return nil
}

// runSQLMigrations executes SQL migration files from the specified directory
func runSQLMigrations(db *gorm.DB, migrationDir string) error {
	// Check if migration directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		// Directory doesn't exist, skip migrations
		return nil
	}

	// Read all SQL files from the migration directory
	var sqlFiles []string
	err := filepath.WalkDir(migrationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".sql") {
			sqlFiles = append(sqlFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migration directory %s: %w", migrationDir, err)
	}

	// Sort files to ensure consistent execution order
	sort.Strings(sqlFiles)

	// Execute each SQL file
	for _, sqlFile := range sqlFiles {
		if err := executeSQLFile(db, sqlFile); err != nil {
			return fmt.Errorf("failed to execute migration file %s: %w", sqlFile, err)
		}
	}

	return nil
}

// executeSQLFile reads and executes a SQL file
func executeSQLFile(db *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", filePath, err)
	}

	sqlContent := string(content)
	if strings.TrimSpace(sqlContent) == "" {
		// Skip empty files
		return nil
	}

	// Execute the SQL content
	if err := db.Exec(sqlContent).Error; err != nil {
		return fmt.Errorf("failed to execute SQL from file %s: %w", filePath, err)
	}

	fmt.Printf("Successfully executed migration: %s\n", filePath)
	return nil
}
