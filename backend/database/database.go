package database

import (
	"hrcs/backend/models"
	"hrcs/backend/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	// Run GORM auto-migrations for basic schema
	err := db.AutoMigrate(
		&models.User{},
		&models.UserGroup{},
		&models.ClaimType{},
		&models.Claim{},
		&models.ApprovalLevel{},
		&models.ClaimApproval{},
		
		// NEW MODELS FOR CLAIMS, BALANCES & DETAILS VIEWS
		&models.UserGroupClaimType{},
		&models.UserClaimType{},
		&models.UserClaimBalance{},
		
		// NEW MODELS FOR PHASE 3 FEATURES
		&models.AuditLog{},
	)
	
	if err != nil {
		return err
	}

	// Run custom SQL migrations for indexes and constraints
	return migrations.RunMigrations(db)
}