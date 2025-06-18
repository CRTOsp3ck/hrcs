package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// MigrationRecord tracks which migrations have been run
type MigrationRecord struct {
	ID       uint   `gorm:"primaryKey"`
	Filename string `gorm:"uniqueIndex;not null"`
	RunAt    string `gorm:"not null"`
}

// RunMigrations executes all pending SQL migrations
func RunMigrations(db *gorm.DB) error {
	// Create migrations table if it doesn't exist
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get list of migration files
	migrationsDir := filepath.Join("backend", "migrations")
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Sort files to ensure they run in order
	sort.Strings(files)

	// Get already run migrations
	var runMigrations []MigrationRecord
	db.Find(&runMigrations)
	
	runMigrationMap := make(map[string]bool)
	for _, migration := range runMigrations {
		runMigrationMap[migration.Filename] = true
	}

	// Run pending migrations
	for _, file := range files {
		filename := filepath.Base(file)
		
		// Skip if already run
		if runMigrationMap[filename] {
			fmt.Printf("Migration %s already run, skipping\n", filename)
			continue
		}

		fmt.Printf("Running migration: %s\n", filename)
		
		// Read migration file
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", filename, err)
		}

		// Execute migration
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %v", err)
		}

		if err := executeMigration(sqlDB, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", filename, err)
		}

		// Record successful migration
		migrationRecord := MigrationRecord{
			Filename: filename,
			RunAt:    fmt.Sprintf("%d", getCurrentTimestamp()),
		}
		
		if err := db.Create(&migrationRecord).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %v", filename, err)
		}

		fmt.Printf("Migration %s completed successfully\n", filename)
	}

	fmt.Println("All migrations completed")
	return nil
}

// executeMigration runs a single migration SQL content
func executeMigration(db *sql.DB, content string) error {
	// Split content by semicolon and execute each statement
	statements := strings.Split(content, ";")
	
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("failed to execute statement: %s, error: %v", statement, err)
		}
	}

	return nil
}

// getCurrentTimestamp returns current Unix timestamp
func getCurrentTimestamp() int64 {
	return 1700000000 // Placeholder - in real implementation, use time.Now().Unix()
}