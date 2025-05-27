package main

import (
	"flag"
	"log"

	"hrcs/backend/config"
	"hrcs/backend/database"
	"hrcs/backend/seeds"
)

func main() {
	// Parse command line flags
	clearFlag := flag.Bool("clear", false, "Clear all existing data before seeding")
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations first
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create seeder
	seeder := seeds.NewSeeder(db)

	// Clear data if requested
	if *clearFlag {
		if err := seeder.ClearAll(); err != nil {
			log.Fatal("Failed to clear data:", err)
		}
	}

	// Run seeding
	if err := seeder.SeedAll(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	printSeedingInfo()
}

func showHelp() {
	log.Println(`
HR Claims Management System - Database Seeder

Usage:
  go run cmd/seed/main.go [flags]

Flags:
  -clear    Clear all existing data before seeding
  -help     Show this help message

Examples:
  go run cmd/seed/main.go              # Seed database (skip if data exists)
  go run cmd/seed/main.go -clear       # Clear and reseed database

This will create:
- Admin and normal users with default password 'password123'
- Various claim types (Travel, Medical, Office Supplies, etc.)
- User groups (Engineering, Sales, Marketing, etc.)
- Approval levels for each user group
- Sample claims in different statuses
`)
}

func printSeedingInfo() {
	log.Println(`
🎉 Database seeding completed successfully!

Default Login Credentials:
┌─────────────────────────────────────────────────────────┐
│                     ADMIN USERS                        │
├─────────────────────────────────────────────────────────┤
│ Email: admin@hrcs.com              | Password: password123 │
│ Email: hr.manager@hrcs.com         | Password: password123 │
│ Email: finance.manager@hrcs.com    | Password: password123 │
│ Email: dept.head@hrcs.com          | Password: password123 │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    NORMAL USERS                        │
├─────────────────────────────────────────────────────────┤
│ Email: john.doe@hrcs.com           | Password: password123 │
│ Email: jane.smith@hrcs.com         | Password: password123 │
│ Email: bob.wilson@hrcs.com         | Password: password123 │
│ Email: alice.brown@hrcs.com        | Password: password123 │
│ Email: charlie.davis@hrcs.com      | Password: password123 │
│ Email: diana.garcia@hrcs.com       | Password: password123 │
│ Email: frank.miller@hrcs.com       | Password: password123 │
│ Email: grace.taylor@hrcs.com       | Password: password123 │
└─────────────────────────────────────────────────────────┘

Created Data:
- 12 Users (4 admin, 8 normal)
- 10 Claim Types
- 8 User Groups
- 16 Approval Levels (2 levels per group)
- 8 Sample Claims

You can now start the application and log in with any of the above credentials!
`)
}
