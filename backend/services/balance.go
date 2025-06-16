package services

import (
	"errors"
	"fmt"
	"hrcs/backend/models"
	"time"

	"gorm.io/gorm"
)

type BalanceService struct {
	db *gorm.DB
}

func NewBalanceService(db *gorm.DB) *BalanceService {
	return &BalanceService{db: db}
}

// GetUserBalance gets or creates a user's balance record for a claim type
func (s *BalanceService) GetUserBalance(userID, claimTypeID uint) (*models.UserClaimBalance, error) {
	var balance models.UserClaimBalance
	
	// Try to find existing balance
	result := s.db.Where("user_id = ? AND claim_type_id = ?", userID, claimTypeID).
		Preload("User").
		Preload("ClaimType").
		First(&balance)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new balance record
			return s.createUserBalance(userID, claimTypeID)
		}
		return nil, result.Error
	}
	
	// Check if balance needs reset
	if balance.NeedsReset() {
		err := s.ResetBalanceIfNeeded(&balance)
		if err != nil {
			return nil, err
		}
	}
	
	return &balance, nil
}

// createUserBalance creates a new balance record for a user and claim type
func (s *BalanceService) createUserBalance(userID, claimTypeID uint) (*models.UserClaimBalance, error) {
	// Get the claim type to determine default limits
	var claimType models.ClaimType
	if err := s.db.First(&claimType, claimTypeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("claim type not found")
		}
		return nil, fmt.Errorf("failed to fetch claim type: %w", err)
	}
	
	// Check for user-specific overrides
	limitAmount := claimType.LimitAmount
	var userOverride models.UserClaimType
	result := s.db.Where("user_id = ? AND claim_type_id = ?", userID, claimTypeID).First(&userOverride)
	if result.Error == nil && userOverride.CustomLimitAmount != nil {
		limitAmount = *userOverride.CustomLimitAmount
	}
	
	// Check for user group overrides
	var user models.User
	if err := s.db.Preload("UserGroup").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	
	if user.UserGroupID != nil {
		var groupOverride models.UserGroupClaimType
		result := s.db.Where("user_group_id = ? AND claim_type_id = ?", *user.UserGroupID, claimTypeID).First(&groupOverride)
		if result.Error == nil && groupOverride.CustomLimitAmount != nil && userOverride.ID == 0 {
			// Only use group override if no user override exists
			limitAmount = *groupOverride.CustomLimitAmount
		}
	}
	
	// Create new balance record
	balance := models.UserClaimBalance{
		UserID:           userID,
		ClaimTypeID:      claimTypeID,
		TotalLimit:       limitAmount,
		CurrentSpent:     0,
		RemainingBalance: limitAmount,
		LastResetDate:    time.Now(),
		ResetPeriod:      claimType.LimitTimespan,
	}
	
	if err := s.db.Create(&balance).Error; err != nil {
		return nil, fmt.Errorf("failed to create balance record: %w", err)
	}
	
	// Load relationships for return
	if err := s.db.Preload("User").Preload("ClaimType").First(&balance, balance.ID).Error; err != nil {
		// Even if we can't load relationships, return the balance we created
		balance.User = user
		balance.ClaimType = claimType
	}
	
	return &balance, nil
}

// CanUserClaim checks if a user can claim the specified amount
func (s *BalanceService) CanUserClaim(userID, claimTypeID uint, amount float64) (bool, string, error) {
	// First check if user has permission to claim this type
	canAccess, err := s.canUserAccessClaimType(userID, claimTypeID)
	if err != nil {
		return false, "", err
	}
	if !canAccess {
		return false, "You do not have permission to claim this type", nil
	}
	
	// Get user's balance
	balance, err := s.GetUserBalance(userID, claimTypeID)
	if err != nil {
		return false, "", err
	}
	
	// Check if amount exceeds remaining balance
	if amount > balance.RemainingBalance {
		return false, fmt.Sprintf("Amount $%.2f exceeds remaining balance of $%.2f", amount, balance.RemainingBalance), nil
	}
	
	return true, "", nil
}

// canUserAccessClaimType checks if user has permission to access a claim type
func (s *BalanceService) canUserAccessClaimType(userID, claimTypeID uint) (bool, error) {
	// Check for user-specific override first
	var userOverride models.UserClaimType
	result := s.db.Where("user_id = ? AND claim_type_id = ?", userID, claimTypeID).First(&userOverride)
	if result.Error == nil {
		return userOverride.IsAllowed, nil
	}
	
	// Check user group permissions
	var user models.User
	if err := s.db.Preload("UserGroup").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("user not found")
		}
		return false, fmt.Errorf("failed to fetch user: %w", err)
	}
	
	if user.UserGroupID != nil {
		var groupPermission models.UserGroupClaimType
		result := s.db.Where("user_group_id = ? AND claim_type_id = ?", *user.UserGroupID, claimTypeID).First(&groupPermission)
		if result.Error == nil {
			return groupPermission.IsAllowed, nil
		}
	}
	
	// Default to allowed if no specific permissions are set
	return true, nil
}

// DeductFromBalance deducts an amount from user's balance (called when claim is paid)
func (s *BalanceService) DeductFromBalance(userID, claimTypeID uint, amount float64) error {
	balance, err := s.GetUserBalance(userID, claimTypeID)
	if err != nil {
		return err
	}
	
	// Update balance
	balance.CurrentSpent += amount
	balance.RemainingBalance = balance.TotalLimit - balance.CurrentSpent
	
	// Ensure remaining balance doesn't go negative
	if balance.RemainingBalance < 0 {
		balance.RemainingBalance = 0
	}
	
	return s.db.Save(balance).Error
}

// ResetBalanceIfNeeded resets balance if the reset period has elapsed
func (s *BalanceService) ResetBalanceIfNeeded(balance *models.UserClaimBalance) error {
	if !balance.NeedsReset() {
		return nil
	}
	
	// Reset the balance
	balance.CurrentSpent = 0
	balance.RemainingBalance = balance.TotalLimit
	balance.LastResetDate = time.Now()
	
	return s.db.Save(balance).Error
}

// GetAllUserBalances gets all balance records for a user
func (s *BalanceService) GetAllUserBalances(userID uint) ([]models.UserClaimBalance, error) {
	var balances []models.UserClaimBalance
	
	err := s.db.Where("user_id = ?", userID).
		Preload("ClaimType").
		Find(&balances).Error
	
	// Reset balances if needed
	for i := range balances {
		if balances[i].NeedsReset() {
			s.ResetBalanceIfNeeded(&balances[i])
		}
	}
	
	return balances, err
}

// AdminAdjustBalance allows admin to manually adjust a user's balance
func (s *BalanceService) AdminAdjustBalance(userID, claimTypeID uint, newLimit float64) error {
	balance, err := s.GetUserBalance(userID, claimTypeID)
	if err != nil {
		return err
	}
	
	// Update the total limit
	balance.TotalLimit = newLimit
	
	// Recalculate remaining balance
	balance.RemainingBalance = balance.TotalLimit - balance.CurrentSpent
	if balance.RemainingBalance < 0 {
		balance.RemainingBalance = 0
	}
	
	return s.db.Save(balance).Error
}