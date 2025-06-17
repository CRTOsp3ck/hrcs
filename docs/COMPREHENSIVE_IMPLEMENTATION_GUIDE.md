# Comprehensive Implementation Guide for High-Priority Features

## Project Analysis Summary

### Current Tech Stack
- **Frontend**: Vue 3 + TypeScript + PrimeVue + TailwindCSS + Vite
- **Backend**: Go + GORM + Gin + PostgreSQL
- **State Management**: Pinia
- **Routing**: Vue Router 4
- **UI Components**: PrimeVue 4.3.4 with custom CSS variables

### Current Architecture Overview

#### Frontend Structure
```
frontend/src/
├── components/         # Reusable components
├── layouts/           # AdminLayout.vue, plus future layouts
├── views/             # Page components
│   ├── admin/         # Admin-specific views
│   └── (user views)   # User-facing views
├── stores/            # Pinia stores (auth.ts, etc.)
├── api/               # API client functions
├── types/             # TypeScript interfaces
└── router/            # Vue Router configuration
```

#### Backend Structure
```
backend/
├── models/            # GORM models (claim.go, user.go)
├── handlers/          # HTTP handlers
├── services/          # Business logic (balance.go)
├── middleware/        # Auth middleware
└── routes/            # Route definitions
```

#### Key Models
- **User**: Has role (admin/normal), belongs to UserGroup
- **Claim**: Status workflow (draft → submitted → approved/rejected → payment-in-progress → paid)
- **ClaimType**: Has limit_amount and limit_timespan for balance tracking
- **UserClaimBalance**: Tracks spending per user per claim type
- **ApprovalLevel**: Multi-level approval workflow system
- **ClaimApproval**: Individual approval records with permissions

---

## IMPLEMENTATION ROADMAP

## 1. NON-ADMIN/NORMAL USER SECTION

### 1.1 Dashboard - Empty State Spacing Fix

**Issue**: Need gap between 'No claims found' and 'Create Your First Claim' button in recent claims empty state.

**Current Location**: `frontend/src/views/DashboardView.vue:91-99`

**Current Code**:
```vue
<template #empty>
  <div class="empty-state">
    <i class="pi pi-inbox empty-icon"></i>
    <p>No claims found</p>
    <Button
      label="Create Your First Claim"
      icon="pi pi-plus"
      @click="router.push('/claims/new')"
    />
  </div>
</template>
```

**Implementation**:
1. **File**: `frontend/src/views/DashboardView.vue`
2. **CSS Update** (line 383-392):
```css
.empty-state {
  text-align: center;
  padding: 3rem;
}

.empty-state p {
  margin-bottom: 1.5rem; /* ADD THIS LINE */
}

.empty-icon {
  font-size: 3rem;
  color: var(--surface-400);
  margin-bottom: 1rem;
}
```

### 1.2 ClaimDetailView - Remove Admin Actions & Add Approval Workflow

**Issues**: 
- Remove admin actions section for non-admin users
- Show approval workflow and history for claim submitters

**Current Location**: `frontend/src/views/ClaimDetailView.vue:126-161`

**Implementation Strategy**:

1. **Conditional Admin Actions Display** (line 126):
```vue
<!-- Admin Actions -->
<Card v-if="authStore.isAdmin && claim.status === 'submitted'" class="admin-actions-card">
```

2. **Enhanced Approval Workflow Section** (replace/enhance lines 96-123):
```vue
<!-- Approval Workflow & History -->
<Card class="workflow-card">
  <template #header>
    <h3 class="card-title">Approval Workflow</h3>
  </template>
  <template #content>
    <!-- Workflow Progress Summary -->
    <div v-if="claim.approvalWorkflow && claim.approvalWorkflow.length > 0" class="workflow-summary">
      <div class="progress-container">
        <ProgressBar 
          :value="getWorkflowProgress()" 
          :showValue="false" 
          style="height: 8px; margin-bottom: 1rem;"
        />
        <div class="progress-text">
          {{ getCompletedSteps() }} of {{ claim.approvalWorkflow.length }} approval steps completed
        </div>
      </div>
      
      <!-- Current Status -->
      <div class="current-status">
        <strong>Current Status:</strong>
        <Tag :value="formatStatus(claim.status)" :severity="getStatusSeverity(claim.status)" />
      </div>
    </div>

    <!-- Detailed Workflow Steps -->
    <Timeline :value="approvalWorkflowSteps" class="workflow-timeline">
      <template #marker="slotProps">
        <span class="workflow-marker" :class="getStepMarkerClass(slotProps.item)">
          <i :class="getStepIcon(slotProps.item)"></i>
        </span>
      </template>
      <template #content="slotProps">
        <div class="workflow-step-content">
          <div class="step-header">
            <h4 class="step-title">{{ slotProps.item.title }}</h4>
            <Tag 
              :value="slotProps.item.statusLabel" 
              :severity="slotProps.item.statusSeverity"
              :icon="slotProps.item.statusIcon"
            />
          </div>
          
          <div class="step-details">
            <p class="step-description">{{ slotProps.item.description }}</p>
            <div class="step-meta">
              <span v-if="slotProps.item.approver" class="approver-info">
                <i class="pi pi-user mr-1"></i>
                {{ slotProps.item.approver }}
              </span>
              <span v-if="slotProps.item.date" class="step-date">
                <i class="pi pi-calendar mr-1"></i>
                {{ formatDateTime(slotProps.item.date) }}
              </span>
            </div>
            <div v-if="slotProps.item.comments" class="step-comments">
              <i class="pi pi-comment mr-1"></i>
              {{ slotProps.item.comments }}
            </div>
          </div>
        </div>
      </template>
    </Timeline>
  </template>
</Card>
```

3. **Add Computed Properties**:
```typescript
const approvalWorkflowSteps = computed(() => {
  if (!claim.value) return []
  
  const steps = [
    {
      id: 'created',
      title: 'Claim Created',
      description: `Claim created by ${claim.value.user?.name}`,
      date: claim.value.created_at,
      status: 'completed',
      statusLabel: 'Created',
      statusSeverity: 'info' as const,
      statusIcon: 'pi pi-plus'
    }
  ]
  
  if (claim.value.submitted_at) {
    steps.push({
      id: 'submitted',
      title: 'Submitted for Approval',
      description: 'Claim submitted to approval workflow',
      date: claim.value.submitted_at,
      status: 'completed',
      statusLabel: 'Submitted',
      statusSeverity: 'info' as const,
      statusIcon: 'pi pi-send'
    })
  }
  
  // Add workflow steps from backend
  if (claim.value.approvalWorkflow) {
    claim.value.approvalWorkflow.forEach((step, index) => {
      steps.push({
        id: `approval-${step.id}`,
        title: `Level ${step.level}: ${step.name}`,
        description: `Assigned to ${step.approverName} (${step.userGroupName})`,
        approver: step.approverName,
        date: step.completedAt,
        status: step.status,
        statusLabel: getStepStatusLabel(step.status),
        statusSeverity: getStepStatusSeverity(step.status),
        statusIcon: getStepStatusIcon(step.status),
        comments: step.comments
      })
    })
  }
  
  return steps
})

const getWorkflowProgress = () => {
  if (!claim.value?.approvalWorkflow) return 0
  const completed = claim.value.approvalWorkflow.filter(step => 
    step.status === 'approved' || step.status === 'rejected'
  ).length
  return (completed / claim.value.approvalWorkflow.length) * 100
}

const getCompletedSteps = () => {
  if (!claim.value?.approvalWorkflow) return 0
  return claim.value.approvalWorkflow.filter(step => 
    step.status === 'approved' || step.status === 'rejected'
  ).length
}
```

### 1.3 NewClaimView - Critical Fixes & Enhancements

**Issues**:
- User gets logged out when selecting claim type
- Claim type dropdown shows 'No available options' when editing
- Missing balance information display
- Need claim limit enforcement

**Current Location**: `frontend/src/views/NewClaimView.vue`

**Implementation Strategy**:

1. **Fix Logout Issue** - Debug claim type loading:
```typescript
// In loadClaimTypes function (line 382-398)
const loadClaimTypes = async () => {
  try {
    const response = await claimTypesApi.getAll()
    console.log('Claim types response:', response.data) // Add debug logging
    
    if (response.data.data) {
      claimTypes.value = response.data.data.filter(type => {
        // Add proper filtering logic
        return type.is_active !== false // Ensure we show active types
      })
      console.log('Loaded claim types:', claimTypes.value) // Debug log
    }
  } catch (error) {
    console.error('Failed to load claim types:', error) // Better error logging
    
    // Check if it's an auth error
    if (error.response?.status === 401) {
      console.log('Auth error when loading claim types')
      // Don't show toast for auth errors - let interceptor handle
      return
    }
    
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claim types',
      life: 3000
    })
  }
}
```

2. **Backend API Endpoint Check** - Ensure `/api/claim-types` endpoint exists and is accessible to normal users.

3. **Enhanced Balance Information** - Current balance card needs styling improvements:
```vue
<!-- Enhanced Balance Information Card (lines 61-91) -->
<Card v-if="selectedClaimType && balanceInfo" class="balance-info-card span-2">
  <template #header>
    <div class="balance-header">
      <i class="pi pi-wallet mr-2"></i>
      <h3>Balance Information - {{ selectedClaimType.name }}</h3>
      <ProgressSpinner v-if="loadingBalance" style="width: 20px; height: 20px" />
    </div>
  </template>
  <template #content>
    <div class="balance-summary">
      <div class="balance-overview">
        <div class="balance-stat">
          <label>Limit Period:</label>
          <Tag :value="selectedClaimType.limit_timespan" severity="info" />
        </div>
        <div class="balance-stat">
          <label>Next Reset:</label>
          <span>{{ getNextResetDate() }}</span>
        </div>
      </div>
    </div>
    
    <div class="balance-grid">
      <div class="balance-item">
        <label>Total Limit ({{ selectedClaimType.limit_timespan }}):</label>
        <span class="amount total-limit">${{ balanceInfo.total_limit.toFixed(2) }}</span>
      </div>
      <div class="balance-item">
        <label>Already Spent:</label>
        <span class="amount spent">${{ balanceInfo.current_spent.toFixed(2) }}</span>
      </div>
      <div class="balance-item">
        <label>Remaining Balance:</label>
        <span class="amount remaining" :class="{ 
          'low-balance': balanceInfo.remaining_balance < (balanceInfo.total_limit * 0.1),
          'zero-balance': balanceInfo.remaining_balance <= 0
        }">
          ${{ balanceInfo.remaining_balance.toFixed(2) }}
        </span>
      </div>
    </div>
    
    <!-- Balance Warnings -->
    <div v-if="balanceInfo.remaining_balance <= 0" class="balance-warning danger">
      <i class="pi pi-exclamation-triangle mr-2"></i>
      <span>No remaining balance for this claim type</span>
    </div>
    <div v-else-if="balanceInfo.remaining_balance < (balanceInfo.total_limit * 0.1)" class="balance-warning warning">
      <i class="pi pi-info-circle mr-2"></i>
      <span>Low balance remaining ({{ ((balanceInfo.remaining_balance / balanceInfo.total_limit) * 100).toFixed(1) }}% left)</span>
    </div>
  </template>
</Card>
```

4. **Enforce Balance Limits** - Update amount input validation:
```vue
<!-- Enhanced Amount Input (lines 42-58) -->
<div class="form-field">
  <label for="amount" class="form-label required">Amount ($)</label>
  <InputNumber
    id="amount"
    v-model="form.amount"
    mode="currency"
    currency="USD"
    locale="en-US"
    :min="0"
    :max="balanceInfo?.remaining_balance || 999999"
    :invalid="!!errors.amount || isBalanceExceeded"
    class="w-full"
    @input="handleAmountChange"
  />
  <small v-if="errors.amount" class="p-error">{{ errors.amount }}</small>
  <small v-if="isBalanceExceeded && !errors.amount" class="p-error">
    Amount exceeds remaining balance of ${{ balanceInfo?.remaining_balance.toFixed(2) }}
  </small>
  <small v-if="balanceInfo && form.amount > 0 && !isBalanceExceeded" class="p-help">
    Remaining after this claim: ${{ (balanceInfo.remaining_balance - form.amount).toFixed(2) }}
  </small>
</div>
```

5. **Add Helper Functions**:
```typescript
const getNextResetDate = () => {
  if (!balanceInfo.value) return 'N/A'
  
  const lastReset = new Date(balanceInfo.value.last_reset_date)
  const resetPeriod = balanceInfo.value.reset_period
  
  let nextReset = new Date(lastReset)
  
  switch (resetPeriod) {
    case 'daily':
      nextReset.setDate(nextReset.getDate() + 1)
      break
    case 'weekly':
      nextReset.setDate(nextReset.getDate() + 7)
      break
    case 'monthly':
      nextReset.setMonth(nextReset.getMonth() + 1)
      break
    case 'annual':
      nextReset.setFullYear(nextReset.getFullYear() + 1)
      break
  }
  
  return nextReset.toLocaleDateString()
}

const handleAmountChange = () => {
  // Real-time balance validation
  if (form.amount && balanceInfo.value) {
    if (form.amount > balanceInfo.value.remaining_balance) {
      form.amount = balanceInfo.value.remaining_balance
      toast.add({
        severity: 'warn',
        summary: 'Amount Adjusted',
        detail: `Amount automatically adjusted to remaining balance of $${balanceInfo.value.remaining_balance.toFixed(2)}`,
        life: 4000
      })
    }
  }
}
```

---

## 2. ADMIN MODE SECTION

### 2.1 General Layout Fixes

**Issue**: Admin layout content needs to shift down due to header proximity

**Current Location**: `frontend/src/layouts/AdminLayout.vue`

**Implementation**:
1. **Update Admin Layout CSS** (line 183-188):
```css
.admin-content {
  background: var(--surface-0);
  border-radius: var(--space-8);
  padding: var(--space-34) var(--space-21) var(--space-21); /* Increase top padding */
  box-shadow: 0 var(--space-1) var(--space-5) rgba(0, 0, 0, 0.1);
}
```

2. **Add Loading State for DataTables** - Update all admin views with tables:
```vue
<!-- In AdminClaims.vue, AdminUsers.vue, etc. -->
<DataTable
  :value="data"
  :loading="loading"
  v-show="!loading || data.length > 0"
  <!-- ... other props -->
>
</DataTable>

<!-- Add loading skeleton -->
<div v-show="loading && data.length === 0" class="table-loading">
  <div class="loading-skeleton">
    <div class="skeleton-row" v-for="i in 5" :key="i">
      <div class="skeleton-cell"></div>
      <div class="skeleton-cell"></div>
      <div class="skeleton-cell"></div>
    </div>
  </div>
</div>
```

### 2.2 Admin Dashboard Enhancements

**Issues**: 
- Move percentage text below icons and figures
- Reorganize to 3-column layout
- Enhance charts with more details
- Add audit log section

**Current Location**: `frontend/src/views/admin/AdminDashboard.vue`

**Implementation**:

1. **Restructure Stats Cards** (lines 9-81):
```vue
<!-- Enhanced Key Metrics -->
<div class="metrics-grid">
  <Card class="metric-card">
    <template #content>
      <div class="metric-content">
        <div class="metric-icon" style="background: rgba(59, 130, 246, 0.1);">
          <i class="pi pi-users" style="color: #3b82f6;"></i>
        </div>
        <div class="metric-details">
          <p class="metric-value">{{ stats.totalUsers }}</p>
          <p class="metric-label">Total Users</p>
          <p class="metric-change positive">
            <i class="pi pi-arrow-up"></i>
            12% from last month
          </p>
        </div>
      </div>
    </template>
  </Card>
  <!-- Repeat for other metrics -->
</div>
```

2. **New 3-Column Dashboard Layout** (replace lines 83-125):
```vue
<div class="dashboard-grid-enhanced">
  <!-- Claims by Status Chart -->
  <Card class="chart-card">
    <template #header>
      <div class="chart-header">
        <h3 class="card-title">Claims by Status</h3>
        <div class="chart-stats">
          <span class="total-count">{{ getTotalClaimsCount() }} total</span>
        </div>
      </div>
    </template>
    <template #content>
      <Chart type="doughnut" :data="claimsByStatusData" :options="enhancedChartOptions" />
      <div class="chart-legend-details">
        <div v-for="item in stats.claimsByStatus" :key="item.status" class="legend-item">
          <div class="legend-color" :style="{ backgroundColor: getStatusColor(item.status) }"></div>
          <span class="legend-label">{{ formatStatus(item.status) }}</span>
          <span class="legend-count">{{ item.count }}</span>
          <span class="legend-percentage">{{ getStatusPercentage(item.count) }}%</span>
        </div>
      </div>
    </template>
  </Card>

  <!-- Claims by Type Chart -->
  <Card class="chart-card">
    <template #header>
      <div class="chart-header">
        <h3 class="card-title">Claims by Type</h3>
        <div class="chart-stats">
          <span class="highest-type">{{ getHighestClaimType() }}</span>
        </div>
      </div>
    </template>
    <template #content>
      <Chart type="bar" :data="claimsByTypeData" :options="enhancedBarChartOptions" />
      <div class="type-summary">
        <div v-for="item in stats.claimsByType" :key="item.type" class="type-item">
          <div class="type-info">
            <span class="type-name">{{ item.type }}</span>
            <span class="type-amount">${{ formatAmount(item.amount) }}</span>
          </div>
          <div class="type-count">{{ item.count }} claims</div>
        </div>
      </div>
    </template>
  </Card>

  <!-- Quick Actions -->
  <Card class="quick-actions-enhanced">
    <template #header>
      <h3 class="card-title">Quick Actions</h3>
    </template>
    <template #content>
      <div class="actions-grid-vertical">
        <Button 
          label="Manage Users" 
          icon="pi pi-users"
          @click="router.push('/admin/users')"
          class="action-button"
        />
        <Button 
          label="Configure Groups" 
          icon="pi pi-sitemap"
          severity="secondary"
          @click="router.push('/admin/groups')"
          class="action-button"
        />
        <Button 
          label="Claim Types" 
          icon="pi pi-tags"
          severity="info"
          @click="router.push('/admin/claim-types')"
          class="action-button"
        />
        <Button 
          label="Approval Workflow" 
          icon="pi pi-shield"
          severity="warning"
          @click="router.push('/admin/approval-levels')"
          class="action-button"
        />
      </div>
    </template>
  </Card>
</div>

<!-- Recent Claims Section -->
<div class="dashboard-row">
  <Card class="recent-claims-admin">
    <template #header>
      <div class="section-header">
        <h3 class="card-title">Recent Claims</h3>
        <Button 
          label="View All Claims" 
          icon="pi pi-arrow-right"
          text
          @click="router.push('/admin/claims')"
        />
      </div>
    </template>
    <template #content>
      <DataTable
        :value="recentClaims"
        :loading="loading"
        responsiveLayout="scroll"
        :paginator="false"
        class="p-datatable-sm recent-claims-table"
      >
        <!-- Table columns here -->
      </DataTable>
    </template>
  </Card>
</div>

<!-- Recent Activity (Audit Log) -->
<div class="dashboard-row">
  <Card class="recent-activity-enhanced">
    <template #header>
      <div class="section-header">
        <h3 class="card-title">Recent Activity</h3>
        <Button 
          label="View Audit Log" 
          icon="pi pi-arrow-right"
          text
          @click="router.push('/admin/audit-log')"
        />
      </div>
    </template>
    <template #content>
      <Timeline :value="recentActivity" class="activity-timeline">
        <template #marker="slotProps">
          <span class="activity-marker" :class="getActivityClass(slotProps.item.type)">
            <i :class="getActivityIcon(slotProps.item.type)"></i>
          </span>
        </template>
        <template #content="slotProps">
          <div class="activity-content">
            <div class="activity-header">
              <p class="activity-title">{{ slotProps.item.title }}</p>
              <span class="activity-user">{{ slotProps.item.user }}</span>
            </div>
            <p class="activity-time">{{ formatTimeAgo(slotProps.item.timestamp) }}</p>
            <p v-if="slotProps.item.details" class="activity-details">{{ slotProps.item.details }}</p>
          </div>
        </template>
      </Timeline>
    </template>
  </Card>
</div>
```

3. **Update CSS** (add to style section):
```css
/* Enhanced Dashboard Layout */
.dashboard-grid-enhanced {
  display: grid;
  grid-template-columns: 1fr 1fr 300px;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.dashboard-row {
  margin-bottom: 2rem;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

/* Enhanced Metric Cards */
.metric-content {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.metric-details {
  flex: 1;
  text-align: left;
}

.metric-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0 0 0.25rem;
  line-height: 1;
}

.metric-label {
  color: var(--surface-600);
  font-size: 0.875rem;
  margin: 0 0 0.5rem;
}

.metric-change {
  font-size: 0.875rem;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* Chart Enhancements */
.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-stats {
  font-size: 0.875rem;
  color: var(--surface-600);
}

.chart-legend-details {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-label {
  flex: 1;
  font-size: 0.875rem;
}

.legend-count {
  font-weight: 600;
  font-size: 0.875rem;
}

.legend-percentage {
  font-size: 0.75rem;
  color: var(--surface-500);
  min-width: 35px;
  text-align: right;
}

/* Quick Actions */
.actions-grid-vertical {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.action-button {
  width: 100%;
  justify-content: flex-start;
}

/* Activity Timeline */
.activity-timeline .activity-marker {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 0.875rem;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.25rem;
}

.activity-user {
  font-size: 0.75rem;
  color: var(--surface-500);
  font-weight: 500;
}

.activity-details {
  font-size: 0.875rem;
  color: var(--surface-600);
  margin: 0.25rem 0 0;
}

@media (max-width: 1400px) {
  .dashboard-grid-enhanced {
    grid-template-columns: 1fr 1fr;
  }
  
  .quick-actions-enhanced {
    grid-column: span 2;
  }
}
```

### 2.3 TODO: Audit Log Feature

**New Feature**: Full audit log implementation

**Implementation Plan**:

1. **Backend Model** (`backend/models/audit.go`):
```go
type AuditLog struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    UserID      uint      `json:"user_id" gorm:"not null"`
    User        User      `json:"user"`
    Action      string    `json:"action" gorm:"not null"`
    EntityType  string    `json:"entity_type" gorm:"not null"`
    EntityID    uint      `json:"entity_id"`
    OldValues   string    `json:"old_values" gorm:"type:jsonb"`
    NewValues   string    `json:"new_values" gorm:"type:jsonb"`
    IPAddress   string    `json:"ip_address"`
    UserAgent   string    `json:"user_agent"`
    CreatedAt   time.Time `json:"created_at"`
}
```

2. **Frontend View** (`frontend/src/views/admin/AdminAuditLog.vue`):
```vue
<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Audit Log</h1>
      <p class="page-subtitle">System activity and change tracking</p>
    </div>

    <!-- Filters -->
    <div class="filters-section">
      <div class="filters-row">
        <Calendar
          v-model="filters.dateRange"
          selectionMode="range"
          dateFormat="mm/dd/yy"
          placeholder="Date Range"
          showIcon
        />
        
        <Dropdown
          v-model="filters.action"
          :options="actionOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="All Actions"
          showClear
        />
        
        <Dropdown
          v-model="filters.entityType"
          :options="entityTypeOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="All Entity Types"
          showClear
        />
        
        <InputText 
          v-model="filters.user"
          placeholder="Search by user..."
        />
      </div>
    </div>

    <!-- Audit Log Table -->
    <DataTable
      :value="auditLogs"
      :loading="loading"
      paginator
      :rows="25"
      :rowsPerPageOptions="[25, 50, 100]"
      responsiveLayout="scroll"
      class="audit-table"
    >
      <Column field="created_at" header="Timestamp" sortable style="width: 150px">
        <template #body="slotProps">
          {{ formatDateTime(slotProps.data.created_at) }}
        </template>
      </Column>
      
      <Column field="user.name" header="User" sortable>
        <template #body="slotProps">
          <div class="user-info">
            <Avatar :label="getInitials(slotProps.data.user.name)" size="small" />
            <span>{{ slotProps.data.user.name }}</span>
          </div>
        </template>
      </Column>
      
      <Column field="action" header="Action" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.action" :severity="getActionSeverity(slotProps.data.action)" />
        </template>
      </Column>
      
      <Column field="entity_type" header="Entity" sortable>
        <template #body="slotProps">
          {{ formatEntityType(slotProps.data.entity_type) }}
        </template>
      </Column>
      
      <Column field="details" header="Details">
        <template #body="slotProps">
          <div class="action-details">
            {{ getActionDescription(slotProps.data) }}
          </div>
        </template>
      </Column>
      
      <Column header="Changes" style="width: 100px">
        <template #body="slotProps">
          <Button
            icon="pi pi-eye"
            text
            rounded
            @click="viewChanges(slotProps.data)"
            v-tooltip="'View Changes'"
          />
        </template>
      </Column>
    </DataTable>
  </div>
</template>
```

### 2.4 TODO: Reports Feature

**New Feature**: Comprehensive reporting system

**Implementation Plan**:

1. **Reports View** (`frontend/src/views/admin/AdminReports.vue`):
```vue
<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Reports</h1>
      <p class="page-subtitle">Generate and view system reports</p>
    </div>

    <!-- Report Types Grid -->
    <div class="reports-grid">
      <Card v-for="report in reportTypes" :key="report.id" class="report-card">
        <template #header>
          <div class="report-header">
            <i :class="report.icon" class="report-icon"></i>
            <h3>{{ report.name }}</h3>
          </div>
        </template>
        <template #content>
          <p class="report-description">{{ report.description }}</p>
          <div class="report-actions">
            <Button 
              label="Generate" 
              @click="generateReport(report.id)"
              :loading="generatingReports.includes(report.id)"
            />
            <Button 
              label="Schedule" 
              severity="secondary"
              @click="scheduleReport(report.id)"
            />
          </div>
        </template>
      </Card>
    </div>

    <!-- Recent Reports -->
    <Card class="recent-reports">
      <template #header>
        <h3 class="card-title">Recent Reports</h3>
      </template>
      <template #content>
        <DataTable :value="recentReports" responsiveLayout="scroll">
          <!-- Report history table -->
        </DataTable>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
const reportTypes = [
  {
    id: 'claims-summary',
    name: 'Claims Summary',
    description: 'Overview of all claims by status, type, and period',
    icon: 'pi pi-chart-pie'
  },
  {
    id: 'user-activity',
    name: 'User Activity',
    description: 'User engagement and claim submission patterns',
    icon: 'pi pi-users'
  },
  {
    id: 'financial-overview',
    name: 'Financial Overview',
    description: 'Total amounts, approvals, and budget utilization',
    icon: 'pi pi-dollar'
  },
  {
    id: 'approval-efficiency',
    name: 'Approval Efficiency',
    description: 'Approval times and workflow performance',
    icon: 'pi pi-clock'
  },
  {
    id: 'compliance-audit',
    name: 'Compliance Audit',
    description: 'Regulatory compliance and audit trail',
    icon: 'pi pi-shield'
  }
]
</script>
```

### 2.5 TODO: Microsoft Integration Placeholders

**New Feature**: SSO and Teams integration preparation

**Implementation Plan**:

1. **Integration Settings View** (`frontend/src/views/admin/AdminIntegrations.vue`):
```vue
<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Integrations</h1>
      <p class="page-subtitle">Configure external system integrations</p>
    </div>

    <!-- Microsoft Integration Section -->
    <Card class="integration-card">
      <template #header>
        <div class="integration-header">
          <i class="pi pi-microsoft" style="color: #0078d4;"></i>
          <div>
            <h3>Microsoft Integration</h3>
            <span class="integration-status">Not Configured</span>
          </div>
        </div>
      </template>
      <template #content>
        <div class="integration-content">
          <div class="feature-list">
            <h4>Available Features:</h4>
            <ul>
              <li>
                <i class="pi pi-check-circle"></i>
                Single Sign-On (SSO) with Microsoft Entra ID
              </li>
              <li>
                <i class="pi pi-check-circle"></i>
                Microsoft Teams notifications
              </li>
              <li>
                <i class="pi pi-check-circle"></i>
                Active Directory user synchronization
              </li>
            </ul>
          </div>
          
          <div class="integration-actions">
            <Message severity="info" :closable="false">
              Microsoft integration is currently in development. Contact your system administrator for setup assistance.
            </Message>
            
            <div class="action-buttons">
              <Button 
                label="Configure SSO" 
                icon="pi pi-cog"
                disabled
                @click="showSSOConfig = true"
              />
              <Button 
                label="Setup Teams" 
                icon="pi pi-comments"
                severity="secondary"
                disabled
                @click="showTeamsConfig = true"
              />
            </div>
          </div>
        </div>
      </template>
    </Card>

    <!-- Other Integrations -->
    <Card class="integration-card">
      <template #header>
        <div class="integration-header">
          <i class="pi pi-cloud" style="color: #6366f1;"></i>
          <div>
            <h3>Other Integrations</h3>
            <span class="integration-status">Coming Soon</span>
          </div>
        </div>
      </template>
      <template #content>
        <div class="coming-soon">
          <p>Additional integrations will be available in future releases:</p>
          <ul>
            <li>Slack notifications</li>
            <li>Email providers (SendGrid, Mailgun)</li>
            <li>Document storage (SharePoint, Google Drive)</li>
            <li>Accounting systems (QuickBooks, SAP)</li>
          </ul>
        </div>
      </template>
    </Card>
  </div>
</template>
```

---

## 3. BACKEND API REQUIREMENTS

### 3.1 New Endpoints Needed

1. **Audit Log API**:
```go
// GET /api/admin/audit-log
// POST /api/audit (internal use for logging)
```

2. **Enhanced Claims API**:
```go
// GET /api/claims/:id/workflow - Get detailed approval workflow
// PUT /api/claims/:id/workflow/:stepId - Update workflow step
```

3. **Reports API**:
```go
// GET /api/admin/reports
// POST /api/admin/reports/:type/generate
// GET /api/admin/reports/:id/download
```

4. **Integration API**:
```go
// GET /api/admin/integrations
// PUT /api/admin/integrations/microsoft
```

### 3.2 Enhanced Response Models

1. **Claim with Workflow**:
```go
type ClaimWithWorkflow struct {
    Claim
    ApprovalWorkflow []ApprovalStep `json:"approval_workflow"`
    CurrentStep      *ApprovalStep  `json:"current_step"`
    NextSteps        []ApprovalStep `json:"next_steps"`
}

type ApprovalStep struct {
    ID              uint                 `json:"id"`
    Level           int                  `json:"level"`
    Name            string               `json:"name"`
    ApproverID      uint                 `json:"approver_id"`
    ApproverName    string               `json:"approver_name"`
    ApproverEmail   string               `json:"approver_email"`
    UserGroupID     uint                 `json:"user_group_id"`
    UserGroupName   string               `json:"user_group_name"`
    Status          ClaimStatus          `json:"status"`
    CompletedAt     *time.Time           `json:"completed_at"`
    Comments        string               `json:"comments"`
    Permissions     ApprovalPermissions  `json:"permissions"`
}
```

---

## 4. IMPLEMENTATION PRIORITY

### Phase 1
1. ✅ Dashboard empty state spacing fix
2. ✅ Remove admin actions from ClaimDetailView for non-admins
3. ✅ Fix NewClaimView logout issue
4. ✅ Admin layout content spacing

### Phase 2
1. Enhanced ClaimDetailView approval workflow
2. NewClaimView balance enforcement
3. Admin dashboard restructure
4. Table loading states

### Phase 3
1. Audit log feature implementation
2. Reports feature framework
3. Microsoft integration placeholders

### Phase 4
1. Real-time notifications
2. Advanced reporting
3. Full Microsoft integration
4. Mobile responsiveness improvements

---

## 5. TESTING STRATEGY

### Frontend Testing
1. **Component Tests**: Test all new components with Vue Test Utils
2. **Integration Tests**: Test API interactions and state management
3. **E2E Tests**: User workflows for claim submission and approval
4. **Accessibility Tests**: Ensure PrimeVue components are accessible

### Backend Testing
1. **Unit Tests**: Test business logic and models
2. **API Tests**: Test all endpoints with various user roles
3. **Database Tests**: Test GORM queries and relationships
4. **Performance Tests**: Load testing for claim processing

### Browser Compatibility
- Chrome 100+
- Firefox 100+
- Safari 15+
- Edge 100+

---

## 6. DEPLOYMENT CONSIDERATIONS

### Environment Variables
```env
# Frontend (.env)
VITE_API_URL=http://localhost:8000/api
VITE_APP_NAME=Claimatic
VITE_MICROSOFT_CLIENT_ID=placeholder

# Backend (.env)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=claimatic
DB_USER=postgres
DB_PASSWORD=password
JWT_SECRET=your-jwt-secret
MICROSOFT_TENANT_ID=placeholder
MICROSOFT_CLIENT_ID=placeholder
MICROSOFT_CLIENT_SECRET=placeholder
```

### Database Migrations
```sql
-- Add audit log table
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add indexes for performance
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
```

---

This comprehensive guide provides detailed implementation instructions for all requested features, prioritized by importance and technical complexity. Each section includes specific file locations, code examples, and implementation strategies to ensure efficient and accurate development.
