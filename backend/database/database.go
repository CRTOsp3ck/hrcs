package database

import (
	"hrcs/backend/models"

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
	return db.AutoMigrate(
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
	)
}