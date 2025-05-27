package seeds

import (
	"fmt"
	"log"

	"hrcs/backend/models"
	"hrcs/backend/utils"

	"gorm.io/gorm"
)

type Seeder struct {
	DB *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{DB: db}
}

func (s *Seeder) SeedAll() error {
	log.Println("🌱 Starting database seeding...")

	if err := s.SeedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := s.SeedClaimTypes(); err != nil {
		return fmt.Errorf("failed to seed claim types: %w", err)
	}

	if err := s.SeedUserGroups(); err != nil {
		return fmt.Errorf("failed to seed user groups: %w", err)
	}

	if err := s.SeedApprovalLevels(); err != nil {
		return fmt.Errorf("failed to seed approval levels: %w", err)
	}

	if err := s.SeedSampleClaims(); err != nil {
		return fmt.Errorf("failed to seed sample claims: %w", err)
	}

	log.Println("✅ Database seeding completed successfully!")
	return nil
}

func (s *Seeder) SeedUsers() error {
	log.Println("👥 Seeding users...")

	// Check if users already exist
	var count int64
	s.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already exist, skipping...")
		return nil
	}

	// Hash password for all users (using 'password123' for demo)
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return err
	}

	users := []models.User{
		// Admin Users
		{
			Email:     "admin@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Super",
			LastName:  "Admin",
			Role:      models.RoleAdmin,
		},
		{
			Email:     "hr.manager@hrcs.com", 
			Password:  hashedPassword,
			FirstName: "Sarah",
			LastName:  "Johnson",
			Role:      models.RoleAdmin,
		},
		{
			Email:     "finance.manager@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Michael",
			LastName:  "Chen",
			Role:      models.RoleAdmin,
		},
		{
			Email:     "dept.head@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Emily",
			LastName:  "Rodriguez",
			Role:      models.RoleAdmin,
		},
		// Normal Users (Employees)
		{
			Email:     "john.doe@hrcs.com",
			Password:  hashedPassword,
			FirstName: "John",
			LastName:  "Doe",
			Role:      models.RoleNormal,
		},
		{
			Email:     "jane.smith@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Jane",
			LastName:  "Smith",
			Role:      models.RoleNormal,
		},
		{
			Email:     "bob.wilson@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Bob",
			LastName:  "Wilson",
			Role:      models.RoleNormal,
		},
		{
			Email:     "alice.brown@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Alice",
			LastName:  "Brown",
			Role:      models.RoleNormal,
		},
		{
			Email:     "charlie.davis@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Charlie",
			LastName:  "Davis",
			Role:      models.RoleNormal,
		},
		{
			Email:     "diana.garcia@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Diana",
			LastName:  "Garcia",
			Role:      models.RoleNormal,
		},
		{
			Email:     "frank.miller@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Frank",
			LastName:  "Miller",
			Role:      models.RoleNormal,
		},
		{
			Email:     "grace.taylor@hrcs.com",
			Password:  hashedPassword,
			FirstName: "Grace",
			LastName:  "Taylor",
			Role:      models.RoleNormal,
		},
	}

	if err := s.DB.Create(&users).Error; err != nil {
		return err
	}

	log.Printf("✅ Created %d users", len(users))
	return nil
}

func (s *Seeder) SeedClaimTypes() error {
	log.Println("📋 Seeding claim types...")

	// Check if claim types already exist
	var count int64
	s.DB.Model(&models.ClaimType{}).Count(&count)
	if count > 0 {
		log.Println("Claim types already exist, skipping...")
		return nil
	}

	claimTypes := []models.ClaimType{
		{
			Name:        "Travel Expenses",
			Description: "Business travel related expenses including flights, hotels, meals, and transportation",
		},
		{
			Name:        "Medical Expenses",
			Description: "Health and medical related expenses covered by company policy",
		},
		{
			Name:        "Office Supplies",
			Description: "Office equipment, stationery, and supplies purchased for work",
		},
		{
			Name:        "Training & Development",
			Description: "Professional development courses, conferences, and training materials",
		},
		{
			Name:        "Entertainment",
			Description: "Client entertainment and business meal expenses",
		},
		{
			Name:        "Technology",
			Description: "Software licenses, hardware, and IT equipment",
		},
		{
			Name:        "Telecommunications",
			Description: "Phone bills, internet, and communication services",
		},
		{
			Name:        "Vehicle Expenses",
			Description: "Fuel, maintenance, and vehicle-related business expenses",
		},
		{
			Name:        "Professional Services",
			Description: "Consulting, legal, and other professional service fees",
		},
		{
			Name:        "Miscellaneous",
			Description: "Other business-related expenses not covered by other categories",
		},
	}

	if err := s.DB.Create(&claimTypes).Error; err != nil {
		return err
	}

	log.Printf("✅ Created %d claim types", len(claimTypes))
	return nil
}

func (s *Seeder) SeedUserGroups() error {
	log.Println("👨‍👩‍👧‍👦 Seeding user groups...")

	// Check if user groups already exist
	var count int64
	s.DB.Model(&models.UserGroup{}).Count(&count)
	if count > 0 {
		log.Println("User groups already exist, skipping...")
		return nil
	}

	userGroups := []models.UserGroup{
		{
			Name:        "Engineering",
			Description: "Software engineers, developers, and technical staff",
		},
		{
			Name:        "Sales",
			Description: "Sales representatives, account managers, and business development",
		},
		{
			Name:        "Marketing",
			Description: "Marketing team, content creators, and digital marketing specialists",
		},
		{
			Name:        "Finance",
			Description: "Financial analysts, accountants, and finance team members",
		},
		{
			Name:        "Human Resources",
			Description: "HR team, recruiters, and people operations",
		},
		{
			Name:        "Operations",
			Description: "Operations team, project managers, and administrative staff",
		},
		{
			Name:        "Management",
			Description: "Senior management, directors, and executive team",
		},
		{
			Name:        "Customer Support",
			Description: "Customer service representatives and support team",
		},
	}

	if err := s.DB.Create(&userGroups).Error; err != nil {
		return err
	}

	// Update users with user groups
	var users []models.User
	if err := s.DB.Find(&users).Error; err != nil {
		return err
	}

	// Assign user groups to normal users
	normalUsers := make([]models.User, 0)
	for _, user := range users {
		if user.Role == models.RoleNormal {
			normalUsers = append(normalUsers, user)
		}
	}

	// Distribute normal users across different groups
	groupAssignments := map[string][]string{
		"Engineering":       {"john.doe@hrcs.com", "bob.wilson@hrcs.com"},
		"Sales":            {"jane.smith@hrcs.com", "charlie.davis@hrcs.com"},
		"Marketing":        {"alice.brown@hrcs.com"},
		"Finance":          {"diana.garcia@hrcs.com"},
		"Operations":       {"frank.miller@hrcs.com"},
		"Customer Support": {"grace.taylor@hrcs.com"},
	}

	for groupName, userEmails := range groupAssignments {
		var group models.UserGroup
		if err := s.DB.Where("name = ?", groupName).First(&group).Error; err != nil {
			continue
		}

		for _, email := range userEmails {
			if err := s.DB.Model(&models.User{}).Where("email = ?", email).Update("user_group_id", group.ID).Error; err != nil {
				log.Printf("Warning: Failed to assign user %s to group %s", email, groupName)
			}
		}
	}

	log.Printf("✅ Created %d user groups and assigned users", len(userGroups))
	return nil
}

func (s *Seeder) SeedApprovalLevels() error {
	log.Println("📝 Seeding approval levels...")

	// Check if approval levels already exist
	var count int64
	s.DB.Model(&models.ApprovalLevel{}).Count(&count)
	if count > 0 {
		log.Println("Approval levels already exist, skipping...")
		return nil
	}

	// Get admin users and user groups
	var adminUsers []models.User
	if err := s.DB.Where("role = ?", models.RoleAdmin).Find(&adminUsers).Error; err != nil {
		return err
	}

	var userGroups []models.UserGroup
	if err := s.DB.Find(&userGroups).Error; err != nil {
		return err
	}

	if len(adminUsers) < 2 || len(userGroups) == 0 {
		log.Println("Not enough admin users or user groups to create approval levels")
		return nil
	}

	// Create approval levels for each user group
	approvalLevels := []models.ApprovalLevel{}

	for _, group := range userGroups {
		// Level 1: Department Head (first admin user)
		if len(adminUsers) > 1 {
			approvalLevels = append(approvalLevels, models.ApprovalLevel{
				Level:       1,
				UserGroupID: group.ID,
				ApproverID:  adminUsers[1].ID, // HR Manager
				CanApprove:  true,
				CanReject:   true,
			})
		}

		// Level 2: Finance Manager (second admin user)
		if len(adminUsers) > 2 {
			approvalLevels = append(approvalLevels, models.ApprovalLevel{
				Level:       2,
				UserGroupID: group.ID,
				ApproverID:  adminUsers[2].ID, // Finance Manager
				CanApprove:  true,
				CanReject:   true,
			})
		}
	}

	if len(approvalLevels) > 0 {
		if err := s.DB.Create(&approvalLevels).Error; err != nil {
			return err
		}
	}

	log.Printf("✅ Created %d approval levels", len(approvalLevels))
	return nil
}

func (s *Seeder) SeedSampleClaims() error {
	log.Println("💰 Seeding sample claims...")

	// Check if claims already exist
	var count int64
	s.DB.Model(&models.Claim{}).Count(&count)
	if count > 0 {
		log.Println("Claims already exist, skipping...")
		return nil
	}

	// Get normal users and claim types
	var normalUsers []models.User
	if err := s.DB.Where("role = ?", models.RoleNormal).Find(&normalUsers).Error; err != nil {
		return err
	}

	var claimTypes []models.ClaimType
	if err := s.DB.Find(&claimTypes).Error; err != nil {
		return err
	}

	if len(normalUsers) == 0 || len(claimTypes) == 0 {
		log.Println("No normal users or claim types found, skipping sample claims")
		return nil
	}

	sampleClaims := []models.Claim{
		{
			Title:       "Business Trip to New York",
			Description: "Travel expenses for client meeting in NYC including flights, hotel, and meals",
			Amount:      1250.00,
			Status:      models.StatusSubmitted,
			UserID:      normalUsers[0].ID,
			ClaimTypeID: claimTypes[0].ID, // Travel Expenses
		},
		{
			Title:       "Annual Health Checkup",
			Description: "Medical expenses for annual health checkup and dental cleaning",
			Amount:      450.00,
			Status:      models.StatusApproved,
			UserID:      normalUsers[1].ID,
			ClaimTypeID: claimTypes[1].ID, // Medical Expenses
		},
		{
			Title:       "Laptop and Accessories",
			Description: "New laptop, external monitor, and keyboard for remote work setup",
			Amount:      2100.00,
			Status:      models.StatusSubmitted,
			UserID:      normalUsers[2].ID,
			ClaimTypeID: claimTypes[5].ID, // Technology
		},
		{
			Title:       "AWS Conference 2024",
			Description: "Registration fee and accommodation for AWS re:Invent conference",
			Amount:      1800.00,
			Status:      models.StatusDraft,
			UserID:      normalUsers[3].ID,
			ClaimTypeID: claimTypes[3].ID, // Training & Development
		},
		{
			Title:       "Client Dinner Meeting",
			Description: "Business dinner with potential client at upscale restaurant",
			Amount:      320.00,
			Status:      models.StatusPaid,
			UserID:      normalUsers[4].ID,
			ClaimTypeID: claimTypes[4].ID, // Entertainment
		},
		{
			Title:       "Office Furniture",
			Description: "Ergonomic chair and standing desk for home office",
			Amount:      750.00,
			Status:      models.StatusRejected,
			UserID:      normalUsers[0].ID,
			ClaimTypeID: claimTypes[2].ID, // Office Supplies
		},
		{
			Title:       "Monthly Phone Bill",
			Description: "Business mobile phone bill for October 2024",
			Amount:      85.00,
			Status:      models.StatusApproved,
			UserID:      normalUsers[1].ID,
			ClaimTypeID: claimTypes[6].ID, // Telecommunications
		},
		{
			Title:       "Gas Receipts",
			Description: "Fuel expenses for client visits during the month",
			Amount:      180.00,
			Status:      models.StatusSubmitted,
			UserID:      normalUsers[2].ID,
			ClaimTypeID: claimTypes[7].ID, // Vehicle Expenses
		},
	}

	if err := s.DB.Create(&sampleClaims).Error; err != nil {
		return err
	}

	log.Printf("✅ Created %d sample claims", len(sampleClaims))
	return nil
}

func (s *Seeder) ClearAll() error {
	log.Println("🧹 Clearing all data...")

	// Delete in reverse order due to foreign key constraints
	tables := []interface{}{
		&models.ClaimApproval{},
		&models.ApprovalLevel{},
		&models.Claim{},
		&models.ClaimType{},
		&models.User{},
		&models.UserGroup{},
	}

	for _, table := range tables {
		if err := s.DB.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
			return err
		}
	}

	log.Println("✅ All data cleared")
	return nil
}