# HRCS Implementation Guide: las vegas

## Overview
This guide provides a comprehensive implementation plan for the three core requirements:
1. **Claims**: Claim type limits with timespan configurations and user group assignments
2. **Balances**: User balance tracking with deduction on payment completion
3. **Details Views**: Admin detail views for users, user groups, and claim types

---

## 1. CLAIMS IMPLEMENTATION

### 1.1 Backend Models Changes

#### Update ClaimType Model (`backend/models/claim.go`)
```go
// Add to ClaimType struct
type ClaimType struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    
    // NEW FIELDS FOR CLAIM LIMITS
    LimitAmount       float64      `json:"limit_amount" gorm:"default:0"`
    LimitTimespan     LimitTimespan `json:"limit_timespan" gorm:"default:'annual'"`
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Add new enum for limit timespan
type LimitTimespan string

const (
    LimitAnnual  LimitTimespan = "annual"
    LimitMonthly LimitTimespan = "monthly"
    LimitWeekly  LimitTimespan = "weekly"
    LimitDaily   LimitTimespan = "daily"
)
```

#### Create UserGroupClaimType Association Model
```go
// New model for user group claim type permissions
type UserGroupClaimType struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    UserGroupID       uint      `json:"user_group_id" gorm:"not null"`
    ClaimTypeID       uint      `json:"claim_type_id" gorm:"not null"`
    IsAllowed         bool      `json:"is_allowed" gorm:"default:true"`
    CustomLimitAmount *float64  `json:"custom_limit_amount"`
    
    UserGroup         UserGroup `json:"user_group" gorm:"foreignKey:UserGroupID"`
    ClaimType         ClaimType `json:"claim_type" gorm:"foreignKey:ClaimTypeID"`
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
```

#### Create UserClaimType Override Model
```go
// New model for individual user overrides
type UserClaimType struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    UserID            uint      `json:"user_id" gorm:"not null"`
    ClaimTypeID       uint      `json:"claim_type_id" gorm:"not null"`
    IsAllowed         bool      `json:"is_allowed" gorm:"default:true"`
    CustomLimitAmount *float64  `json:"custom_limit_amount"`
    
    User              User      `json:"user" gorm:"foreignKey:UserID"`
    ClaimType         ClaimType `json:"claim_type" gorm:"foreignKey:ClaimTypeID"`
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
```

### 1.2 Backend API Changes

#### Update Admin Handler (`backend/handlers/admin.go`)
```go
// Add new endpoints for claim type management
func (h *AdminHandler) UpdateClaimTypeWithLimits(w http.ResponseWriter, r *http.Request) {
    // Implementation for updating claim type with limit configurations
}

func (h *AdminHandler) SetUserGroupClaimPermissions(w http.ResponseWriter, r *http.Request) {
    // Implementation for setting which claim types a user group can access
}

func (h *AdminHandler) SetUserClaimOverrides(w http.ResponseWriter, r *http.Request) {
    // Implementation for individual user claim type overrides
}
```

#### Update Routes (`backend/routes/routes.go`)
```go
// Add new routes to admin section
adminRouter.PUT("/claim-types/{id}/limits", adminHandler.UpdateClaimTypeWithLimits)
adminRouter.POST("/user-groups/{id}/claim-permissions", adminHandler.SetUserGroupClaimPermissions)
adminRouter.POST("/users/{id}/claim-overrides", adminHandler.SetUserClaimOverrides)
```

### 1.3 Frontend Changes

#### Update ClaimType Interface (`frontend/src/types/index.ts`)
```typescript
export interface ClaimType {
  id: number;
  name: string;
  description?: string;
  active: boolean;
  
  // NEW FIELDS
  limit_amount: number;
  limit_timespan: 'annual' | 'monthly' | 'weekly' | 'daily';
}

export interface UserGroupClaimPermission {
  id: number;
  user_group_id: number;
  claim_type_id: number;
  is_allowed: boolean;
  custom_limit_amount?: number;
}
```

#### Update AdminClaimTypes.vue
- Add limit_amount and limit_timespan fields to the form
- Add section for managing user group permissions
- Add individual user override management

---

## 2. BALANCES IMPLEMENTATION

### 2.1 Backend Models Changes

#### Update User Model (`backend/models/user.go`)
```go
// Add to User struct
type User struct {
    ID           uint       `json:"id" gorm:"primaryKey"`
    Email        string     `json:"email" gorm:"uniqueIndex;not null"`
    Password     string     `json:"-" gorm:"not null"`
    FirstName    string     `json:"first_name" gorm:"not null"`
    LastName     string     `json:"last_name" gorm:"not null"`
    Role         UserRole   `json:"role" gorm:"default:normal"`
    UserGroupID  *uint      `json:"user_group_id"`
    UserGroup    *UserGroup `json:"user_group,omitempty" gorm:"omitempty"`
    
    // NEW FIELDS FOR BALANCE TRACKING
    ClaimBalances []UserClaimBalance `json:"claim_balances" gorm:"foreignKey:UserID"`
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
```

#### Create UserClaimBalance Model
```go
// New model for tracking user balances per claim type
type UserClaimBalance struct {
    ID              uint      `json:"id" gorm:"primaryKey"`
    UserID          uint      `json:"user_id" gorm:"not null"`
    ClaimTypeID     uint      `json:"claim_type_id" gorm:"not null"`
    
    // Balance tracking
    TotalLimit      float64   `json:"total_limit" gorm:"not null"`
    CurrentSpent    float64   `json:"current_spent" gorm:"default:0"`
    RemainingBalance float64  `json:"remaining_balance" gorm:"default:0"`
    
    // Reset tracking
    LastResetDate   time.Time `json:"last_reset_date"`
    ResetPeriod     LimitTimespan `json:"reset_period" gorm:"not null"`
    
    // Relationships
    User            User      `json:"user" gorm:"foreignKey:UserID"`
    ClaimType       ClaimType `json:"claim_type" gorm:"foreignKey:ClaimTypeID"`
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Helper method to check if balance needs reset
func (ucb *UserClaimBalance) NeedsReset() bool {
    now := time.Now()
    switch ucb.ResetPeriod {
    case LimitDaily:
        return !isSameDay(ucb.LastResetDate, now)
    case LimitWeekly:
        return !isSameWeek(ucb.LastResetDate, now)
    case LimitMonthly:
        return !isSameMonth(ucb.LastResetDate, now)
    case LimitAnnual:
        return !isSameYear(ucb.LastResetDate, now)
    default:
        return false
    }
}
```

### 2.2 Backend Services

#### Create Balance Service (`backend/services/balance.go`)
```go
type BalanceService struct {
    db *gorm.DB
}

func NewBalanceService(db *gorm.DB) *BalanceService {
    return &BalanceService{db: db}
}

// Get user's current balance for a claim type
func (s *BalanceService) GetUserBalance(userID, claimTypeID uint) (*UserClaimBalance, error) {
    // Implementation to get or create user balance record
}

// Check if user can claim the specified amount
func (s *BalanceService) CanUserClaim(userID, claimTypeID uint, amount float64) (bool, string, error) {
    // Implementation for balance validation
}

// Deduct amount from user's balance (called when claim is paid)
func (s *BalanceService) DeductFromBalance(userID, claimTypeID uint, amount float64) error {
    // Implementation for balance deduction
}

// Reset balance if needed based on timespan
func (s *BalanceService) ResetBalanceIfNeeded(balance *UserClaimBalance) error {
    // Implementation for automatic balance reset
}
```

### 2.3 Backend API Changes

#### Update Claim Handler (`backend/handlers/claim.go`)
```go
// Modify CreateClaim to include balance validation
func (h *ClaimHandler) CreateClaim(w http.ResponseWriter, r *http.Request) {
    // Existing validation...
    
    // NEW: Balance validation
    canClaim, message, err := h.balanceService.CanUserClaim(userID, req.ClaimTypeID, req.Amount)
    if err != nil {
        utils.WriteErrorResponse(w, http.StatusInternalServerError, "Balance check failed")
        return
    }
    if !canClaim {
        utils.WriteErrorResponse(w, http.StatusBadRequest, message)
        return
    }
    
    // Continue with claim creation...
}

// Add balance deduction on payment completion
func (h *ClaimHandler) MarkClaimAsPaid(w http.ResponseWriter, r *http.Request) {
    // Existing claim update logic...
    
    // NEW: Deduct from balance when marking as paid
    if claim.Status == "paid" {
        err = h.balanceService.DeductFromBalance(claim.UserID, claim.ClaimTypeID, claim.Amount)
        if err != nil {
            utils.WriteErrorResponse(w, http.StatusInternalServerError, "Balance deduction failed")
            return
        }
    }
    
    // Continue...
}
```

#### Add Balance Handler (`backend/handlers/balance.go`)
```go
type BalanceHandler struct {
    balanceService *services.BalanceService
}

func (h *BalanceHandler) GetUserBalances(w http.ResponseWriter, r *http.Request) {
    // Get all balances for a user
}

func (h *BalanceHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
    // Get balance for specific claim type
}

func (h *BalanceHandler) AdminUpdateBalance(w http.ResponseWriter, r *http.Request) {
    // Admin endpoint to adjust user balances
}
```

### 2.4 Frontend Changes

#### Update Types (`frontend/src/types/index.ts`)
```typescript
export interface UserClaimBalance {
  id: number;
  user_id: number;
  claim_type_id: number;
  total_limit: number;
  current_spent: number;
  remaining_balance: number;
  last_reset_date: string;
  reset_period: 'annual' | 'monthly' | 'weekly' | 'daily';
  claim_type: ClaimType;
}

export interface BalanceInfo {
  claim_type_id: number;
  claim_type_name: string;
  total_limit: number;
  current_spent: number;
  remaining_balance: number;
  can_claim_amount: number;
}
```

#### Update NewClaimView.vue
```vue
<!-- Add balance display section -->
<div class="balance-info" v-if="selectedClaimType">
  <Card>
    <template #title>Balance Information</template>
    <template #content>
      <div class="balance-details">
        <div class="balance-item">
          <label>Total Limit:</label>
          <span class="balance-amount">${{ balanceInfo.total_limit.toFixed(2) }}</span>
        </div>
        <div class="balance-item">
          <label>Current Spent:</label>
          <span class="balance-amount">${{ balanceInfo.current_spent.toFixed(2) }}</span>
        </div>
        <div class="balance-item">
          <label>Remaining Balance:</label>
          <span class="balance-amount balance-remaining">
            ${{ balanceInfo.remaining_balance.toFixed(2) }}
          </span>
        </div>
      </div>
    </template>
  </Card>
</div>

<!-- Update amount input with validation -->
<FormField label="Amount" required>
  <InputNumber
    v-model="form.amount"
    mode="currency"
    currency="USD"
    locale="en-US"
    :max="balanceInfo.remaining_balance"
    :class="{ 'p-invalid': amountExceedsBalance }"
  />
  <small class="p-error" v-if="amountExceedsBalance">
    Amount cannot exceed remaining balance of ${{ balanceInfo.remaining_balance.toFixed(2) }}
  </small>
</FormField>
```

#### Update API calls (`frontend/src/api/index.ts`)
```typescript
// Balance-related API calls
export const getUserBalances = (): Promise<ApiResponse<UserClaimBalance[]>> =>
  api.get('/balances');

export const getUserBalance = (claimTypeId: number): Promise<ApiResponse<UserClaimBalance>> =>
  api.get(`/balances/claim-type/${claimTypeId}`);
```

---

## 3. DETAILS VIEWS IMPLEMENTATION

### 3.1 UserDetailsView Component

#### Create UserDetailsView.vue (`frontend/src/views/admin/UserDetailsView.vue`)
```vue
<template>
  <div class="user-details-view">
    <PageHeader>
      <template #title>User Details: {{ user?.first_name }} {{ user?.last_name }}</template>
      <template #actions>
        <Button @click="editUser" icon="pi pi-pencil" label="Edit User" />
        <Button @click="resetPassword" icon="pi pi-key" label="Reset Password" severity="secondary" />
      </template>
    </PageHeader>

    <div class="details-grid">
      <!-- Basic Information -->
      <Card class="basic-info">
        <template #title>Basic Information</template>
        <template #content>
          <div class="info-grid">
            <div class="info-item">
              <label>Full Name:</label>
              <span>{{ user?.first_name }} {{ user?.last_name }}</span>
            </div>
            <div class="info-item">
              <label>Email:</label>
              <span>{{ user?.email }}</span>
            </div>
            <div class="info-item">
              <label>Role:</label>
              <Tag :value="user?.role" :severity="getRoleSeverity(user?.role)" />
            </div>
            <div class="info-item">
              <label>User Group:</label>
              <span>{{ user?.user_group?.name || 'No group assigned' }}</span>
            </div>
            <div class="info-item">
              <label>Member Since:</label>
              <span>{{ formatDate(user?.created_at) }}</span>
            </div>
          </div>
        </template>
      </Card>

      <!-- Claim Balances -->
      <Card class="balance-info">
        <template #title>Claim Balances</template>
        <template #content>
          <DataTable :value="userBalances" responsiveLayout="scroll">
            <Column field="claim_type.name" header="Claim Type" />
            <Column field="total_limit" header="Total Limit">
              <template #body="slotProps">
                ${{ slotProps.data.total_limit.toFixed(2) }}
              </template>
            </Column>
            <Column field="current_spent" header="Current Spent">
              <template #body="slotProps">
                ${{ slotProps.data.current_spent.toFixed(2) }}
              </template>
            </Column>
            <Column field="remaining_balance" header="Remaining">
              <template #body="slotProps">
                <span :class="getBalanceClass(slotProps.data.remaining_balance)">
                  ${{ slotProps.data.remaining_balance.toFixed(2) }}
                </span>
              </template>
            </Column>
            <Column field="reset_period" header="Reset Period" />
          </DataTable>
        </template>
      </Card>

      <!-- Claims History -->
      <Card class="claims-history">
        <template #title>Recent Claims</template>
        <template #content>
          <DataTable :value="userClaims" responsiveLayout="scroll">
            <Column field="title" header="Title" />
            <Column field="claim_type.name" header="Type" />
            <Column field="amount" header="Amount">
              <template #body="slotProps">
                ${{ slotProps.data.amount.toFixed(2) }}
              </template>
            </Column>
            <Column field="status" header="Status">
              <template #body="slotProps">
                <Tag :value="slotProps.data.status" :severity="getStatusSeverity(slotProps.data.status)" />
              </template>
            </Column>
            <Column field="created_at" header="Created">
              <template #body="slotProps">
                {{ formatDate(slotProps.data.created_at) }}
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <!-- Permissions Override -->
      <Card class="permissions-override">
        <template #title>Claim Type Permissions</template>
        <template #content>
          <div class="permissions-list">
            <div v-for="permission in userPermissions" :key="permission.claim_type_id" class="permission-item">
              <div class="permission-info">
                <span class="claim-type-name">{{ permission.claim_type.name }}</span>
                <Tag :value="permission.is_allowed ? 'Allowed' : 'Denied'" 
                     :severity="permission.is_allowed ? 'success' : 'danger'" />
              </div>
              <div class="custom-limit" v-if="permission.custom_limit_amount">
                Custom Limit: ${{ permission.custom_limit_amount.toFixed(2) }}
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
// Component implementation with data fetching and methods
</script>
```

### 3.2 UserGroupDetailsView Component

#### Create UserGroupDetailsView.vue (`frontend/src/views/admin/UserGroupDetailsView.vue`)
```vue
<template>
  <div class="user-group-details-view">
    <PageHeader>
      <template #title>User Group Details: {{ userGroup?.name }}</template>
      <template #actions>
        <Button @click="editGroup" icon="pi pi-pencil" label="Edit Group" />
        <Button @click="manageMembers" icon="pi pi-users" label="Manage Members" severity="secondary" />
      </template>
    </PageHeader>

    <div class="details-grid">
      <!-- Basic Information -->
      <Card class="basic-info">
        <template #title>Basic Information</template>
        <template #content>
          <div class="info-grid">
            <div class="info-item">
              <label>Group Name:</label>
              <span>{{ userGroup?.name }}</span>
            </div>
            <div class="info-item">
              <label>Description:</label>
              <span>{{ userGroup?.description || 'No description' }}</span>
            </div>
            <div class="info-item">
              <label>Total Members:</label>
              <span>{{ groupMembers?.length || 0 }}</span>
            </div>
            <div class="info-item">
              <label>Created:</label>
              <span>{{ formatDate(userGroup?.created_at) }}</span>
            </div>
          </div>
        </template>
      </Card>

      <!-- Group Members -->
      <Card class="group-members">
        <template #title>Group Members</template>
        <template #content>
          <DataTable :value="groupMembers" responsiveLayout="scroll">
            <Column field="first_name" header="Name">
              <template #body="slotProps">
                {{ slotProps.data.first_name }} {{ slotProps.data.last_name }}
              </template>
            </Column>
            <Column field="email" header="Email" />
            <Column field="role" header="Role">
              <template #body="slotProps">
                <Tag :value="slotProps.data.role" :severity="getRoleSeverity(slotProps.data.role)" />
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button @click="viewUserDetails(slotProps.data.id)" 
                        icon="pi pi-eye" text size="small" />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <!-- Claim Type Permissions -->
      <Card class="claim-permissions">
        <template #title>Claim Type Permissions</template>
        <template #content>
          <DataTable :value="groupClaimPermissions" responsiveLayout="scroll">
            <Column field="claim_type.name" header="Claim Type" />
            <Column field="is_allowed" header="Allowed">
              <template #body="slotProps">
                <Tag :value="slotProps.data.is_allowed ? 'Yes' : 'No'" 
                     :severity="slotProps.data.is_allowed ? 'success' : 'danger'" />
              </template>
            </Column>
            <Column field="custom_limit_amount" header="Custom Limit">
              <template #body="slotProps">
                {{ slotProps.data.custom_limit_amount ? 
                   '$' + slotProps.data.custom_limit_amount.toFixed(2) : 
                   'Default' }}
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button @click="editPermission(slotProps.data)" 
                        icon="pi pi-pencil" text size="small" />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <!-- Approval Levels -->
      <Card class="approval-levels">
        <template #title>Approval Levels</template>
        <template #content>
          <DataTable :value="groupApprovalLevels" responsiveLayout="scroll">
            <Column field="level" header="Level" />
            <Column field="approver.first_name" header="Approver">
              <template #body="slotProps">
                {{ slotProps.data.approver.first_name }} {{ slotProps.data.approver.last_name }}
              </template>
            </Column>
            <Column field="permissions" header="Permissions">
              <template #body="slotProps">
                <div class="permission-tags">
                  <Tag v-if="slotProps.data.can_approve" value="Approve" severity="success" />
                  <Tag v-if="slotProps.data.can_reject" value="Reject" severity="danger" />
                  <Tag v-if="slotProps.data.can_set_paid" value="Set Paid" severity="info" />
                </div>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>
  </div>
</template>
```

### 3.3 ClaimTypeDetailsView Component

#### Create ClaimTypeDetailsView.vue (`frontend/src/views/admin/ClaimTypeDetailsView.vue`)
```vue
<template>
  <div class="claim-type-details-view">
    <PageHeader>
      <template #title>Claim Type Details: {{ claimType?.name }}</template>
      <template #actions>
        <Button @click="editClaimType" icon="pi pi-pencil" label="Edit Claim Type" />
        <Button @click="toggleActive" 
                :icon="claimType?.active ? 'pi pi-eye-slash' : 'pi pi-eye'"
                :label="claimType?.active ? 'Deactivate' : 'Activate'"
                :severity="claimType?.active ? 'danger' : 'success'" />
      </template>
    </PageHeader>

    <div class="details-grid">
      <!-- Basic Information -->
      <Card class="basic-info">
        <template #title>Basic Information</template>
        <template #content>
          <div class="info-grid">
            <div class="info-item">
              <label>Name:</label>
              <span>{{ claimType?.name }}</span>
            </div>
            <div class="info-item">
              <label>Description:</label>
              <span>{{ claimType?.description || 'No description' }}</span>
            </div>
            <div class="info-item">
              <label>Status:</label>
              <Tag :value="claimType?.active ? 'Active' : 'Inactive'" 
                   :severity="claimType?.active ? 'success' : 'danger'" />
            </div>
            <div class="info-item">
              <label>Limit Amount:</label>
              <span>${{ claimType?.limit_amount?.toFixed(2) || '0.00' }}</span>
            </div>
            <div class="info-item">
              <label>Limit Timespan:</label>
              <Tag :value="claimType?.limit_timespan" severity="info" />
            </div>
            <div class="info-item">
              <label>Created:</label>
              <span>{{ formatDate(claimType?.created_at) }}</span>
            </div>
          </div>
        </template>
      </Card>

      <!-- Statistics -->
      <Card class="statistics">
        <template #title>Statistics</template>
        <template #content>
          <div class="stats-grid">
            <StatCard title="Total Claims" :value="stats.total_claims" icon="pi pi-file" />
            <StatCard title="Approved Claims" :value="stats.approved_claims" icon="pi pi-check" />
            <StatCard title="Total Amount" :value="'$' + stats.total_amount.toFixed(2)" icon="pi pi-dollar" />
            <StatCard title="Average Amount" :value="'$' + stats.average_amount.toFixed(2)" icon="pi pi-chart-bar" />
          </div>
        </template>
      </Card>

      <!-- User Group Permissions -->
      <Card class="group-permissions">
        <template #title>User Group Permissions</template>
        <template #content>
          <DataTable :value="groupPermissions" responsiveLayout="scroll">
            <Column field="user_group.name" header="User Group" />
            <Column field="is_allowed" header="Allowed">
              <template #body="slotProps">
                <Tag :value="slotProps.data.is_allowed ? 'Yes' : 'No'" 
                     :severity="slotProps.data.is_allowed ? 'success' : 'danger'" />
              </template>
            </Column>
            <Column field="custom_limit_amount" header="Custom Limit">
              <template #body="slotProps">
                {{ slotProps.data.custom_limit_amount ? 
                   '$' + slotProps.data.custom_limit_amount.toFixed(2) : 
                   'Default' }}
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button @click="editGroupPermission(slotProps.data)" 
                        icon="pi pi-pencil" text size="small" />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>

      <!-- Recent Claims -->
      <Card class="recent-claims">
        <template #title>Recent Claims</template>
        <template #content>
          <DataTable :value="recentClaims" responsiveLayout="scroll">
            <Column field="title" header="Title" />
            <Column field="user.first_name" header="User">
              <template #body="slotProps">
                {{ slotProps.data.user.first_name }} {{ slotProps.data.user.last_name }}
              </template>
            </Column>
            <Column field="amount" header="Amount">
              <template #body="slotProps">
                ${{ slotProps.data.amount.toFixed(2) }}
              </template>
            </Column>
            <Column field="status" header="Status">
              <template #body="slotProps">
                <Tag :value="slotProps.data.status" :severity="getStatusSeverity(slotProps.data.status)" />
              </template>
            </Column>
            <Column field="created_at" header="Created">
              <template #body="slotProps">
                {{ formatDate(slotProps.data.created_at) }}
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>
  </div>
</template>
```

### 3.4 Router Updates

#### Update Router (`frontend/src/router/index.ts`)
```typescript
// Add new detail view routes
{
  path: '/admin/users/:id',
  name: 'UserDetails',
  component: () => import('@/views/admin/UserDetailsView.vue'),
  meta: { requiresAuth: true, requiresAdmin: true }
},
{
  path: '/admin/groups/:id',
  name: 'UserGroupDetails',
  component: () => import('@/views/admin/UserGroupDetailsView.vue'),
  meta: { requiresAuth: true, requiresAdmin: true }
},
{
  path: '/admin/claim-types/:id',
  name: 'ClaimTypeDetails',
  component: () => import('@/views/admin/ClaimTypeDetailsView.vue'),
  meta: { requiresAuth: true, requiresAdmin: true }
}
```

---

## 4. DATABASE AUTO-MIGRATION

### 4.1 GORM Auto-Migration Setup

The database schema will be automatically updated when GORM's auto-migrate runs with the new models. Ensure the following models are registered in your migration function:

#### Update `backend/database/database.go` or migration function:
```go
// Add new models to auto-migration
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.UserGroup{},
        &models.ClaimType{},
        &models.Claim{},
        &models.ApprovalLevel{},
        &models.ClaimApproval{},
        
        // NEW MODELS - Add these for auto-migration
        &models.UserGroupClaimType{},
        &models.UserClaimType{},
        &models.UserClaimBalance{},
    )
}
```

### 4.2 Model Index and Constraint Setup

GORM will automatically create the following based on the model definitions:
- **Foreign key constraints** from `gorm:"foreignKey:..."` tags
- **Unique constraints** from `gorm:"uniqueIndex"` tags  
- **Indexes** from `gorm:"index"` tags
- **Default values** from `gorm:"default:..."` tags

No manual SQL migration scripts are needed - GORM handles all schema changes automatically.

---

## 5. IMPLEMENTATION STEPS

### Phase 1: Backend Foundation
1. Update models with new fields and relationships
2. Add new models to GORM auto-migration
3. Create balance service for calculations and validations
4. Add new API endpoints for balance management
5. Update existing claim creation/approval endpoints with balance checks

### Phase 2: Admin Management
1. Implement admin endpoints for claim type limit configuration
2. Add user group permission management
3. Create individual user override functionality
4. Test all admin functionality

### Phase 3: Balance Integration
1. Update claim creation to show balance information
2. Add balance validation to prevent over-claiming
3. Implement balance deduction on payment completion
4. Add balance reset functionality based on timespan

### Phase 4: Detail Views
1. Create UserDetailsView component
2. Create UserGroupDetailsView component  
3. Create ClaimTypeDetailsView component
4. Update router and navigation
5. Add links from existing admin pages

### Phase 5: Testing & Polish
1. Comprehensive testing of all new features
2. UI/UX refinements
3. Error handling improvements  
4. Documentation updates

---

## 6. API ENDPOINTS SUMMARY

### New Endpoints to Implement
```
Backend Endpoints:
GET    /api/balances                    - Get user's claim balances
GET    /api/balances/claim-type/:id     - Get balance for specific claim type
PUT    /admin/claim-types/:id/limits    - Update claim type limits
POST   /admin/user-groups/:id/claim-permissions - Set group claim permissions
POST   /admin/users/:id/claim-overrides - Set individual user overrides
GET    /admin/users/:id/details         - Get detailed user information
GET    /admin/groups/:id/details        - Get detailed group information
GET    /admin/claim-types/:id/details   - Get detailed claim type information
POST   /admin/balances/adjust           - Admin balance adjustment
```

---

## 7. CRITICAL SUCCESS FACTORS

1. **Balance Deduction Timing**: Ensure balance is only deducted when claim status changes to "paid"
2. **Validation Hierarchy**: User overrides should take precedence over group permissions
3. **Reset Logic**: Implement proper date calculations for different timespans
4. **Concurrent Access**: Handle concurrent claim submissions with proper locking
5. **Data Integrity**: Ensure all balance calculations are atomic and consistent
6. **Performance**: Optimize balance queries for large user bases
7. **Error Handling**: Provide clear error messages when limits are exceeded

This comprehensive guide provides a detailed roadmap for implementing all three core requirements with proper backend architecture, frontend integration, and database design.

<function_calls>
<invoke name="TodoWrite">
<parameter name="todos">[{"id": "1", "content": "Analyze project structure and identify key files", "status": "completed", "priority": "high"}, {"id": "2", "content": "Examine backend models and API structure", "status": "completed", "priority": "high"}, {"id": "3", "content": "Examine frontend components and routing", "status": "completed", "priority": "high"}, {"id": "4", "content": "Analyze claims functionality", "status": "completed", "priority": "high"}, {"id": "5", "content": "Analyze user and user group management", "status": "completed", "priority": "high"}, {"id": "6", "content": "Create comprehensive implementation guide", "status": "completed", "priority": "high"}]
